// Package netrc provides functions to read .netrc files.
//
// netrc syntax consists of groups of pairs of a noun and an identifier
//
//	machine <name>
//	login <name>
//	password <name>
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
//	username, password := netrc.ReadConfig(endpoint, ".netrc", netrc.DefaultNetRC)
package netrc

import (
	"bufio"
	"bytes"
	"io"
	urlpkg "net/url"
	"os"
	"path/filepath"
)

const NetRC = ".netrc"

// DefaultNetRC is "~/.netrc"
var DefaultNetRC = filepath.Join(os.Getenv("HOME"), NetRC)

// ReadConfig reads login and password configuration from file(s), typically ~/.netrc.
// The uri specifies the target system hostname. This can be specified as an endpoint URL or
// simply as a host or domain name.
//
// The netrc files are searched in sequence to find the first match.
//
// The returned username (i.e. login) and password are blank unless a match was found.
func ReadConfig(uri, netrc1 string, netrc ...string) (username string, password string) {
	return readConfig(func(f string) (io.ReadCloser, error) { return os.Open(f) }, uri, netrc1, netrc...)
}

func readConfig(open func(string) (io.ReadCloser, error), uri, netrc1 string, netrc ...string) (string, string) {
	host := uri
	url, err := urlpkg.Parse(uri)
	if err == nil {
		host = url.Host
	}

	if len(netrc) == 0 {
		netrc = []string{}
	}

	var file io.ReadCloser
	file, err = open(netrc1)
	if err == nil {
		u, p, ok := parseConfig(file, host)
		file.Close()
		if ok {
			return u, p
		}
	}

	for _, name := range netrc {
		file, err = open(name)
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
	scanner := bufio.NewScanner(&dropComments{inner: file})
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

type dropComments struct {
	inner   io.Reader
	comment bool
}

func (d *dropComments) Read(p []byte) (n int, err error) {
	n, err = d.inner.Read(p)
	if err != nil {
		return n, err
	}

	hash := 0
	for {
		if d.comment {
			nl := bytes.IndexByte(p[hash:n], '\n') + hash
			if nl < 0 {
				return n, err
			}
			d.comment = false
			copy(p[hash:], p[nl+1:])
			n += hash - (nl + 1)
			hash = 0
		} else {
			hash = bytes.IndexByte(p[:n], '#')
			if hash < 0 {
				return n, err
			}
			d.comment = true
		}
	}
}
