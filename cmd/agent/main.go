package main

import (
	"github.com/XiovV/dokkup/pkg/config"
	"github.com/XiovV/dokkup/pkg/docker"
	"github.com/XiovV/dokkup/pkg/runner"
	"github.com/XiovV/dokkup/pkg/server"
	"go.uber.org/zap"
)

const version = "v0.1.0-beta"

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	hashedKey, err := config.CheckAPIKey()
	if err != nil {
		logger.Fatal("could not check API key", zap.Error(err))
	}

	config, err := config.ReadAgentConfig()
	if err != nil {
		logger.Fatal("could not read config", zap.Error(err))
	}

	config.APIKey = hashedKey

	controller, err := docker.NewController(logger)
	if err != nil {
		logger.Fatal("could not instantiate controller", zap.Error(err))
	}

	jobRunner := runner.NewJobRunner(controller, logger)

	srv := server.Server{Config: config, Controller: controller, JobRunner: jobRunner, Logger: logger}

	logger.Info("server is listening...", zap.String("port", config.Port), zap.String("version", version))
	if err := srv.Serve(); err != nil {
		logger.Fatal("could not start server", zap.Error(err))
	}
}
