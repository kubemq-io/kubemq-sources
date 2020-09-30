package builder

type Manifest struct {
	Schema     string       `json:"schema"`
	Version    string       `json:"version"`
	Connectors []*Connector `json:"connectors"`
}

func (m *Manifest) SetSchema(value string) *Manifest {
	m.Schema = value
	return m
}

func (m *Manifest) SetVersion(value string) *Manifest {
	m.Version = value
	return m
}

func (m *Manifest) AddConnector(value *Connector) *Manifest {
	m.Connectors = append(m.Connectors, value)
	return m
}
