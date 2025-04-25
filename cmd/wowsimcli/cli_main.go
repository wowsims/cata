package main

import (
	"github.com/wowsims/mop/cmd/wowsimcli/cmd"
	"github.com/wowsims/mop/sim"
)

func init() {
	sim.RegisterAll()
}

// Version information.
// This variable is set by the makefile in the release process.
var Version string

func main() {
	cmd.Execute(Version)
}
