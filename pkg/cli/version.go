package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func (a *App) getVersion(ctx *cli.Context) error {
	fmt.Printf("Dokkup v%s, build %s %s\n", Version, Commit[:7], Date)
	return nil
}
