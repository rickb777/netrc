package netrc

import (
	"io"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/rickb777/expect"
)

func TestParseConfig(t *testing.T) {
	u, _ := url.Parse("https://my.server.com:444/foo")
	cases := map[string]string{
		// blank case
		"|": ``,

		// simple one-liner
		"alpha|secret": `machine my.server.com:444 login alpha password secret`,

		// ignore default and match machine
		"beta|secret": `			
			machine my.server.com
			  login xyz
			  password xyz123
			
			machine my.server.com:444
			  login beta 
			  password secret
			  account acct

			default
			  login a1
			  password aaa111`,

		// ignore machine and match default
		"gamma|secret": `			
			machine other.server.com login xyz password xyz123
			default login gamma password secret`,
	}
	for e, input := range cases {
		l, p, _ := parseConfig(strings.NewReader(input), u.Host)
		exp := strings.Split(e, "|")
		expect.String(l).ToBe(t, exp[0])
		expect.String(p).ToBe(t, exp[1])
	}
}

func TestReadConfig(t *testing.T) {
	open := func(s string) (io.ReadCloser, error) {
		if s != "netrc" {
			return nil, os.ErrNotExist
		}
		return io.NopCloser(strings.NewReader(`machine my.server.com:444 login alpha password secret`)), nil
	}

	u, p := readConfig(open, "my.server.com:444", "netrc")
	expect.String(u).ToBe(t, "alpha")
	expect.String(p).ToBe(t, "secret")

	u, p = readConfig(open, "my.server.com:444", "foo", "netrc")
	expect.String(u).ToBe(t, "alpha")
	expect.String(p).ToBe(t, "secret")

	u, p = readConfig(open, "my.server.com:444", "foo", "bar")
	expect.String(u).ToBe(t, "")
	expect.String(p).ToBe(t, "")
}
