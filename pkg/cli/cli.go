package cli

import (
	"bufio"
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
		Usage: "manage containers",
		Commands: []*cli.Command{
			{
				Name:    "run",
				Aliases: []string{"j"},
				Subcommands: []*cli.Command{
					{
						Flags:  defaultFlags(),
						Name:   "job",
						Usage:  "run a job",
						Action: app.jobCmd,
					},
				},
			},
			{
				Name:    "stop",
				Aliases: []string{"s"},
				Subcommands: []*cli.Command{
					{
						Flags:  defaultFlags(),
						Name:   "job",
						Usage:  "stop a job",
						Action: app.stopJobCmd,
					},
				},
			},
			{
				Name:    "rollback",
				Aliases: []string{"r"},
				Subcommands: []*cli.Command{
					{
						Flags:  defaultFlags(),
						Name:   "job",
						Usage:  "rollback a job",
						Action: app.rollbackCmd,
					},
				},
			},
		},
	}

	return app
}

func defaultFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "inventory",
			Aliases: []string{"i"},
			Value:   "inventory.yaml",
			Usage:   "name for the inventory file",
		},
		&cli.BoolFlag{
			Name:    "yes",
			Aliases: []string{"y"},
			Value:   false,
		},
	}
}

func (a *App) Run() {
	if err := a.Cli.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func (a *App) showConfirmationPrompt(ctx *cli.Context) (bool, error) {
	if !ctx.Bool("yes") {
		fmt.Print("\nAre you sure you want to proceed? (y/n) ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return false, err
		}

		if input == "n\n" {
			return false, nil
		}
	}

	return true, nil
}
