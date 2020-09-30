package builder

type Connector struct {
	Kind        string      `json:"kind"`
	Description string      `json:"description"`
	Properties  []*Property `json:"properties"`
}

func NewConnector() *Connector {
	return &Connector{}
}

func (c *Connector) SetKind(value string) *Connector {
	c.Kind = value
	return c
}

func (c *Connector) SetDescription(value string) *Connector {
	c.Description = value
	return c
}
func (c *Connector) AddProperty(value *Property) *Connector {
	c.Properties = append(c.Properties, value)
	return c
}
