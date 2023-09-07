package main

import (
	"github.com/XiovV/dokkup/pkg/config"
	"github.com/XiovV/dokkup/pkg/docker"
	"github.com/XiovV/dokkup/pkg/runner"
	"github.com/XiovV/dokkup/pkg/server"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	key, err := config.CheckAPIKey()
	if err != nil {
		logger.Fatal("could not check API key", zap.Error(err))
	}

	config, err := config.ReadAgentConfig()
	if err != nil {
		logger.Fatal("could not read config", zap.Error(err))
	}

	config.APIKey = key

	controller, err := docker.NewController(logger)
	if err != nil {
		logger.Fatal("could not instantiate controller", zap.Error(err))
	}

	jobRunner := runner.NewJobRunner(controller)

	srv := server.Server{Config: config, Controller: controller, JobRunner: jobRunner, Logger: logger}

	logger.Info("server is listening...", zap.String("port", config.Port))
	if err := srv.Serve(); err != nil {
		logger.Fatal("could not start server", zap.Error(err))
	}
}
