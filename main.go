package main

import (
	_ "embed"
	"github.com/vczyh/dbbackup/cmd"
)

//go:embed VERSION
var version string

func main() {
	cmd.Execute(version)
}
