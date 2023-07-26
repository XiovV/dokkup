package app

import (
	"fmt"
	"log"
	"os"

	"github.com/XiovV/dokkup/config"
	"github.com/urfave/cli/v2"
)

type App struct {
  Cli *cli.App
  Inventory *config.Inventory
} 

func NewApp(inventory *config.Inventory) *App {
  app := &App{Inventory: inventory}

	app.Cli = &cli.App{
		Name:  "dokkup",
		Usage: "test usage",
		Action: func(*cli.Context) error {
			fmt.Println("hello there")
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "add",
				Aliases: []string{"a"},
				Action: app.addCmd,
      },
		},
	}

  return app
}

func (a *App) Run() {
  if err := a.Cli.Run(os.Args); err != nil {
    log.Fatal(err)
  }
}
