package main

import "github.com/synchro-food/filelint/cli"

//go:generate go-bindata -pkg config -o config/bindata.go config/default.yml

func main() {
	cli.Execute()
}
