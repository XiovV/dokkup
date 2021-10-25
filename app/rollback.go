package app

import (
	"fmt"
	"github.com/XiovV/dokkup/controller"
	"os"
)

type Rollback struct {
	node       string
	container  string
	controller controller.DockerController
}

func NewRollback(node, container string, dockerController controller.DockerController) Rollback {
	return Rollback{node: node, container: container, controller: dockerController}
}

func (r *Rollback) Run() {
	fmt.Println("rolling back container...")
	err := r.controller.RollbackContainer(r.container)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("successfully rolled back container")
}
