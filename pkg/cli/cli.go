package cli

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/XiovV/dokkup/pkg/config"
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
						Flags: []cli.Flag{
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
							&cli.BoolFlag{
								Name:    "purge",
								Aliases: []string{"p"},
								Value:   false,
							},
						},
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

func (a *App) readJobAndInventory(ctx *cli.Context) (*config.Job, *config.Inventory, error) {
	inventoryFlag := ctx.String("inventory")
	inventory, err := config.ReadInventory(inventoryFlag)
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't read inventory: %w", err)
	}

	jobArg := ctx.Args().First()

	if len(jobArg) == 0 {
		return nil, nil, errors.New("please provide a job file")
	}

	job, err := config.ReadJob(jobArg)
	if err != nil {
		return nil, nil, err
	}

	return job, inventory, nil
}
