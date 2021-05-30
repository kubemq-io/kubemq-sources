// +build !container

package main

import (
	"github.com/kardianos/service"
	"github.com/kubemq-hub/kubemq-sources/pkg/logger"
	"os"
)

type appService struct {
	serviceExit chan bool
	errCh       chan error
	svclogger   service.Logger
}

func newAppService() *appService {
	a := &appService{
		serviceExit: make(chan bool, 1),
		errCh:       make(chan error, 5),
	}
	return a
}
func (a *appService) init(action, username, password string) error {
	options := make(service.KeyValue)
	options["Restart"] = "on-success"
	options["SuccessExitStatus"] = "1 2 8 SIGKILL"
	options["Password"] = password
	svcConfig := &service.Config{
		Name:             "KubeMQ Sources",
		DisplayName:      "KubeMQ Sources",
		Description:      "KubeMQ Sources connects external systems and cloud services with KubeMQ message queue broker",
		UserName:         username,
		Arguments:        nil,
		Executable:       "",
		Dependencies:     []string{},
		WorkingDirectory: "",
		ChRoot:           "",
		Option:           options,
	}
	srv, err := service.New(a, svcConfig)
	if err != nil {
		return err
	}
	a.svclogger, err = srv.Logger(a.errCh)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			err := <-a.errCh
			if err != nil {
				log.Error(err)
			}
		}
	}()

	if action != "" {
		err := service.Control(srv, action)
		if err != nil {
			return err
		}
		log.Infof("service command %s completed", action)
		return nil
	}
	err = srv.Run()
	return err
}
func (a *appService) Start(s service.Service) error {
	if service.Interactive() {
		preRun()
		log.Infof("starting kubemq sources connectors version: %s", version)
	} else {
		logger.ServiceLog.SetLogger(a.svclogger)
		log.Infof("starting kubemq sources connectors as a service version: %s", version)
	}
	go func() {
		err := runInteractive(a.serviceExit)
		if err != nil {
			log.Errorf("kubemq sources ended with error: %s", err)
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}()
	return nil
}

func (a *appService) Stop(s service.Service) error {
	a.serviceExit <- true
	log.Infof("kubemq sources connectors stopped")
	return nil
}
