package app

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func (a *App) addCmd(ctx *cli.Context) error {
  fmt.Println("hey there new structure")
  return nil
}
