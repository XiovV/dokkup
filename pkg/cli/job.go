package cli

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/XiovV/dokkup/pkg/config"
	"github.com/urfave/cli/v2"
)

func (a *App) jobCmd(ctx *cli.Context) error {
  inventoryFlag := ctx.String("inventory")
  inventory, err := config.ReadInventory(inventoryFlag)
  if err != nil {
    log.Fatal("couldn't read inventory:", err)
  }

  jobArg := ctx.Args().First()

  if len(jobArg) == 0 {
    log.Fatal("please provide a job file")
  }

	job, err := config.ReadJob(jobArg)
	if err != nil {
		return err
	}

  fmt.Println("Deployment summary:\n")

  jobSummaryTable := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)
  fmt.Fprintln(jobSummaryTable, "NAME\tIMAGE\tRESTART\tCOUNT\tGROUP\tNETWORKS")

  out := fmt.Sprintf("%s\t%s\t%s\t%d\t%s\t%s\n", job.Container[0].Name, job.Container[0].Image, job.Container[0].Restart, job.Count, job.Group, job.Container[0].Networks)
  fmt.Fprintln(jobSummaryTable, out)

  jobSummaryTable.Flush()
  
  if (!ctx.Bool("yes")) {
    fmt.Print("Are you sure you want to proceed? (y/n) ")
    reader := bufio.NewReader(os.Stdin)
    input, err := reader.ReadString('\n')
    if err != nil {
      log.Fatal("couldn't read input: ", err)
    }

    if (input == "n\n") {
      return nil
    }
  }

	return nil
}
