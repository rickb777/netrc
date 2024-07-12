// package netrc provides functions to read .netrc files.
//
// netrc syntax consists of pairs of words
//
//	"machine" name
//	"login" name
//	"password" name
//
// (others are ignored here)
//
// The separating whitespace can optionally include newlines.
// The order of the nouns is normally "machine" then "login" then
// "password", but we allow "login" and "password" to be swapped.
//
// Example .netrc content:
//
//	machine foo.com login user@example.com password a123456
//
// Example usage:
//
//	endpoint := "http://my.server.com/
//	username, password := netrc.ReadConfig(endpoint, ".netrc", os.Getenv("HOME")+"/.netrc")
package netrc

import (
	"bufio"
	"io"
	urlpkg "net/url"
	"os"
)

// ReadConfig reads login and password configuration from file(s), typically ~/.netrc.
// The uri specifies the target system, which can be specified as an endpoint URL or
// simply as a host/domain name.
// The returned username and password are blank unless a match was found.
func ReadConfig(uri string, netrc ...string) (string, string) {
	host := uri
	url, err := urlpkg.Parse(uri)
	if err == nil {
		host = url.Host
	}

	for _, name := range netrc {
		var file io.ReadCloser
		file, err = os.Open(name)
		if err == nil {
			u, p, ok := parseConfig(file, host)
			file.Close()
			if ok {
				return u, p
			}
		}
	}

	return "", ""
}

func parseConfig(file io.Reader, host string) (string, string, bool) {
	var machine, login, password string
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	inDefault := false

	for scanner.Scan() {
		id := ""
		noun := scanner.Text()
		if noun == "default" {
			inDefault = true
		} else if scanner.Scan() {
			id = scanner.Text()
		}

		switch noun {
		case "machine":
			inDefault = false
			if id == host {
				machine = id
				login = ""
				password = ""
			}
		case "login":
			if inDefault || machine == host {
				login = id
			}
		case "password":
			if inDefault || machine == host {
				password = id
			}
		}

		if machine == host && login != "" && password != "" {
			return login, password, true
		}
	}

	return login, password, false
}
