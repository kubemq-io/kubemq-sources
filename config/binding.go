package config

import "fmt"

type Binding struct {
	Source string
	Target string
}

func (b Binding) Validate() error {
	if b.Source == "" {
		return fmt.Errorf("binding source cannot be empty")
	}
	if b.Target == "" {
		return fmt.Errorf("binding target cannot be empty")
	}
	return nil
}
