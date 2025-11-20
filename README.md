# netrc

[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg)](http://pkg.go.dev/github.com/rickb777/netrc)
[![Go Report Card](https://goreportcard.com/badge/github.com/rickb777/netrc)](https://goreportcard.com/report/github.com/rickb777/netrc)
[![Build](https://github.com/rickb777/netrc/actions/workflows/go.yml/badge.svg)](https://github.com/rickb777/netrc/actions)
[![Coverage](https://coveralls.io/repos/github/rickb777/netrc/badge.svg?branch=main)](https://coveralls.io/github/rickb777/netrc?branch=main)
[![Issues](https://img.shields.io/github/issues/rickb777/netrc.svg)](https://github.com/rickb777/netrc/issues)

A small API to read `.netrc` file using Go.

# What is .netrc?

A `.netrc` file is a plain-text configuration file used by command-line programs like ftp and curl to store login
credentials for automated connections to remote servers. It often resides in the user's home directory and holds a list
of entries, each specifying a hostname, login, and password, to facilitate passwordless logins for file transfers and
other network operations.

## How it works

- Location: The file is typically stored in the user's home directory, named `.netrc`
- Format: Each entry consists of one or more lines, with the keywords `machine`, `login`, and `password` followed by the
  corresponding values.
- Purpose: When a program like `ftp` or `curl` needs to connect to a remote host, it searches the `.netrc` file for an entry
  that matches the hostname. If a match is found, the stored login and password are used automatically, eliminating the
  need to type them manually.
- Security: The file should have strict permissions to prevent unauthorized access, such as read access denied for the
  group and others. Typically, the file permissions should be `600` (only the owner can read/write the file).
- Example: A simple entry might look like this:

```
machine myhost.com
login myuser
password mysecretpassword
```

This allows the program to log in to `myhost.com` with the username `myuser` and the password `mysecretpassword` without a
prompt.

