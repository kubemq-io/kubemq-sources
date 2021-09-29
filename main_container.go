//go:build container
// +build container

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/kubemq-io/kubemq-sources/api"
	"github.com/kubemq-io/kubemq-sources/binding"
	"github.com/kubemq-io/kubemq-sources/config"
	"github.com/kubemq-io/kubemq-sources/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

var (
	version = ""
)

var (
	log        *logger.Logger
	configFile = flag.String("config", "config.yaml", "set config file name")
)

func run() error {
	var gracefulShutdown = make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGTERM)
	signal.Notify(gracefulShutdown, syscall.SIGINT)
	signal.Notify(gracefulShutdown, syscall.SIGQUIT)
	configCh := make(chan *config.Config)
	cfg, err := config.Load(configCh)
	if err != nil {
		return err
	}
	err = cfg.Validate()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	bindingsService, err := binding.New()
	if err != nil {
		return err
	}
	err = bindingsService.Start(ctx, cfg)
	if err != nil {
		return err
	}

	apiServer, err := api.Start(ctx, cfg.ApiPort, bindingsService)
	if err != nil {
		return err
	}
	for {
		select {
		case newConfig := <-configCh:
			err = newConfig.Validate()
			if err != nil {
				return fmt.Errorf("error on validation new config file: %s", err.Error())

			}
			bindingsService.Stop()
			err = bindingsService.Start(ctx, newConfig)
			if err != nil {
				return fmt.Errorf("error on restarting service with new config file: %s", err.Error())
			}
			if apiServer != nil {
				err = apiServer.Stop()
				if err != nil {
					return fmt.Errorf("error on shutdown api server: %s", err.Error())
				}
			}

			apiServer, err = api.Start(ctx, newConfig.ApiPort, bindingsService)
			if err != nil {
				return fmt.Errorf("error on start api server: %s", err.Error())
			}
		case <-gracefulShutdown:
			_ = apiServer.Stop()
			bindingsService.Stop()
			return nil
		}
	}
}

func main() {
	log = logger.NewLogger("kubemq-sources")
	flag.Parse()
	config.SetConfigFile(*configFile)
	log.Infof("starting kubemq sources connectors version: %s", version)
	if err := run(); err != nil {
		log.Error(err)
		os.Exit(1)
	}

}
