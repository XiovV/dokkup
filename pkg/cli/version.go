package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var (
	Version = "dev"
	Commit  = "none"
)

func (a *App) getVersion(ctx *cli.Context) error {
	fmt.Println("temp println")
	fmt.Printf("Dokkup v%s, build %s\n", Version, Commit[:7])
	return nil
}
