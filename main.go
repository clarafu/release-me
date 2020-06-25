package main

import (
	"github.com/clarafu/release-me/cmd"
)

// Grab all the PRs created after commit of last tag
// Compile all the PR titles and authors into one output
// Output the result

func main() {
	cmd.Execute()
}
