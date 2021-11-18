package main

import (
	"fmt"
	"github.com/XiovV/dokkup/app"
	"github.com/XiovV/dokkup/controller"
	"google.golang.org/grpc"
	"log"
	"os"
)



func main() {
	config := app.NewConfig()
	dockerController := controller.NewDockerController(config.NodeLocation, config.APIKey)

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial("localhost:8080", opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	switch os.Args[1] {
	case "update":
		app := app.NewUpdate(config, dockerController, conn)
		app.Run()
	case "rollback":
		app := app.NewRollback(config, dockerController, conn)
		app.Run()
	default:
		fmt.Println("unknown command")
		os.Exit(1)
	}
}
