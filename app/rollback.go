package app

import (
	"fmt"
	"github.com/XiovV/dokkup/controller"
	"os"
)

type Rollback struct {
	config     *Config
	controller controller.DockerController
}

func NewRollback(config *Config, dockerController controller.DockerController) Rollback {
	return Rollback{config: config, controller: dockerController}
}

func (r *Rollback) Run() {
	fmt.Println("rolling back container...")
	err := r.controller.RollbackContainer(r.config.Container)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("successfully rolled back container")
}
