package app

import (
	"fmt"
	"log"

	"github.com/XiovV/dokkup/config"
	"github.com/urfave/cli/v2"
)

func (a *App) jobCmd(ctx *cli.Context) error {
  job, err := config.ReadJob("../../job.yaml")
  if err != nil {
    log.Fatal(err)
  }
  
  fmt.Println(job)
  
  return nil
}
