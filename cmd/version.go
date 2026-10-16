package cmd

import "fmt"

// Version is the FlakeHunter version, overridable at build time via
// -ldflags "-X github.com/gitstq/flakehunter/cmd.Version=vX.Y.Z".
var Version = "dev"

// VersionCmd implements `flakehunter version`.
func VersionCmd(args []string) error {
	fmt.Printf("flakehunter %s\n", Version)
	fmt.Println("https://github.com/gitstq/flakehunter")
	return nil
}
