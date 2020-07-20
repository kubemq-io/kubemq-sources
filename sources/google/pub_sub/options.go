package google

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/config"
)


type options struct {
	projectID string
	subID string
	concurrency     int
}

func parseOptions(cfg config.Metadata) (options, error) {
	o := options{}
	var err error
	o.projectID , err = cfg.MustParseString("project_id")
	if err != nil {
		return options{}, fmt.Errorf("error parsing project_id, %w", err)
	}
	o.subID , err = cfg.MustParseString("sub_id")
	if err != nil {
		return options{}, fmt.Errorf("error parsing project_id, %w", err)
	}
	err = config.MustExistsEnv("GOOGLE_APPLICATION_CREDENTIALS")
	if err !=nil{
		return options{}, err
	}
	o.concurrency, err = cfg.MustParseIntWithRange("concurrency", 1, 100)
	if err != nil {
		return options{}, fmt.Errorf("error parsing concurrency value, %w", err)
	}
	
	return o, nil
}
