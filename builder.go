package main

import (
	"github.com/kubemq-hub/builder/common"
	"github.com/kubemq-hub/builder/connector"
	"github.com/kubemq-hub/kubemq-sources/sources"
	"github.com/kubemq-hub/kubemq-sources/targets"
	"io/ioutil"
)

func saveManifest() error {
	return common.NewManifest().
		SetSchema("sources").
		SetVersion(version).
		SetSourceConnectors(sources.Connectors()).
		SetTargetConnectors(targets.Connectors()).
		Save("manifest.json")
}

func buildConfig() error {
	if err := saveManifest(); err != nil {
		return err
	}
	var err error
	var bindingsYaml []byte
	if bindingsYaml, err = connector.NewSource().SetManifestFile("./manifest.json").Render(); err != nil {
		return err
	}
	return ioutil.WriteFile("config.yaml", bindingsYaml, 0644)
}
