package main

// VERSION is stamped at build time via -ldflags "-X main.VERSION=...".
// Defaults to a dev marker for plain `go build`.
var VERSION = "v1.2.1-dev"
