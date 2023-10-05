package cli

import (
	"log"

	"github.com/urfave/cli/v2"
)

func (a *App) showJobCmd(ctx *cli.Context) error {
	job, inventory, err := a.readJobAndInventory(ctx)
	if err != nil {
		log.Fatal(err)
	}

	jobStatuses, err := a.getJobStatuses(job, inventory)
	if err != nil {
		log.Fatal(err)
	}

	a.showJobInfoTable(jobStatuses, job)

	return nil
}
