package netrc_test

import (
	"fmt"

	"github.com/rickb777/netrc"
)

func ExampleReadConfig_home_directory() {
	endpoint := "http://my.server.com/"
	// This searches only in the home directory.
	username, password := netrc.ReadConfig(endpoint, netrc.DefaultNetRC)
	fmt.Printf("username=%s\n", username)
	fmt.Printf("password=%s\n", password)
}

func ExampleReadConfig_two_files() {
	endpoint := "http://my.server.com/"
	// This searches two files:
	//  1. in the current directory, and then
	//  2. in the home directory (but only if the first doesn't exist or doesn't contain the endpoint).
	username, password := netrc.ReadConfig(endpoint, "./.netrc", netrc.DefaultNetRC)
	fmt.Printf("username=%s\n", username)
	fmt.Printf("password=%s\n", password)
}
