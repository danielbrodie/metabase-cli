package main

import "github.com/danielbrodie/metabase-cli/cmd"

var version = "dev"

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}
