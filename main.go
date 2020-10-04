package main

import (
	"context"
	"flag"
	"github.com/kubemq-hub/builder/common"
	"github.com/kubemq-hub/kubemq-sources/api"
	"github.com/kubemq-hub/kubemq-sources/binding"
	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/kubemq-hub/kubemq-sources/pkg/logger"
	"github.com/kubemq-hub/kubemq-sources/sources"
	"github.com/kubemq-hub/kubemq-sources/targets"
	"os"
	"os/signal"
	"syscall"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var (
	log              *logger.Logger
	generateManifest = flag.Bool("manifest", false, "generate source connectors manifest")
)

func saveManifest() error {
	return common.NewManifest().
		SetSchema("sources").
		SetVersion(version).
		SetSourceConnectors(sources.Connectors()).
		SetTargetConnectors(targets.Connectors()).
		Save("manifest.json")
}

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
			err = cfg.Validate()
			if err != nil {
				log.Errorf("error on validation new config file: %s", err.Error())
				continue
			}
			bindingsService.Stop()
			err = bindingsService.Start(ctx, newConfig)
			if err != nil {
				log.Errorf("error on restarting service with new config file: %s", err.Error())
				continue
			}
			if apiServer != nil {
				err = apiServer.Stop()
				if err != nil {
					log.Errorf("error on shutdown api server: %s", err.Error())
					continue
				}
			}

			apiServer, err = api.Start(ctx, newConfig.ApiPort, bindingsService)
			if err != nil {
				log.Errorf("error on start api server: %s", err.Error())
				continue
			}
		case <-gracefulShutdown:
			_ = apiServer.Stop()
			bindingsService.Stop()
			return nil
		}
	}
}
func main() {
	log = logger.NewLogger("main")
	flag.Parse()
	if *generateManifest {

		err := saveManifest()
		if err != nil {
			log.Error(err)
			os.Exit(1)

		}
		log.Infof("generated manifest.json completed")
		os.Exit(0)
	}
	log = logger.NewLogger("main")
	log.Infof("starting kubemq sources connectors version: %s, commit: %s, date %s", version, commit, date)
	if err := run(); err != nil {
		log.Error(err)
		os.Exit(1)
	}

}
