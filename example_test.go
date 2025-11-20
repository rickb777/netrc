package netrc_test

import (
	"fmt"
	"os"

	"github.com/rickb777/netrc"
)

func ExampleReadConfig_home_directory() {
	endpoint := "http://my.server.com/"
	// This searches only in the home directory.
	username, password := netrc.ReadConfig(endpoint, netrc.DefaultNetRC)
	fmt.Printf("username=%s\n", username)
	fmt.Printf("password=%s\n", password)
}

func ExampleReadConfig_three_files() {
	endpoint := "http://my.server.com/"

	// This searches up to three files. It ignores absent files and returns as soon as the first match is
	// found, if any.
	username, password := netrc.ReadConfig(endpoint,
		"./.netrc",         // in the current directory
		os.Getenv("NETRC"), // the file location given by NETRC, which is the conventional environment variable to use
		netrc.DefaultNetRC) // i.e. "~/.netrc" in the home directory

	fmt.Printf("username=%s\n", username)
	fmt.Printf("password=%s\n", password)
}
