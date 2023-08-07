package cli

import (
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
		Usage: "manage containers",
		Commands: []*cli.Command{
			{
				Name:    "run",
				Aliases: []string{"j"},
				Subcommands: []*cli.Command{
					{
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "inventory",
                Aliases: []string{"i"},
								Value: "inventory.yaml",
								Usage: "name for the inventory file",
							},
              &cli.BoolFlag{
                Name: "yes",
                Aliases: []string{"y"},
                Value: false,
              },
						},

						Name:   "job",
						Usage:  "run a job",
						Action: app.jobCmd,
					},
				},
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
