package main

import "github.com/fohte/filelint/cli"

//go:generate go-bindata -pkg config -o config/bindata.go config/default.yml

func main() {
	cli.Execute()
}
