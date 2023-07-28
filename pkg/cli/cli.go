package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

type App struct {
  Cli *cli.App
} 

func NewApp() *App {
  app := &App{}

	app.Cli = &cli.App{
		Name:  "dokkup",
		Usage: "test usage",
		Action: func(*cli.Context) error {
			fmt.Println("hello there")
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "job",
				Aliases: []string{"j"},
				Action: app.jobCmd,
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