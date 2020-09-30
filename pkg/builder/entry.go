package builder

type Property struct {
	Name        string   `json:"name"`
	Kind        string   `json:"kind"`
	Description string   `json:"description"`
	Default     string   `json:"default"`
	Options     []string `json:"options"`
	Must        bool     `json:"must"`
	Min         int      `json:"min"`
	Max         int      `json:"max"`
}

func NewProperty() *Property {
	return &Property{}
}

func (p *Property) SetName(value string) *Property {
	p.Name = value
	return p
}

func (p *Property) SetKind(value string) *Property {
	p.Kind = value
	return p
}
func (p *Property) SetDescription(value string) *Property {
	p.Description = value
	return p
}

func (p *Property) SetDefault(value string) *Property {
	p.Default = value
	return p
}

func (p *Property) SetOptions(value string) *Property {
	p.Default = value
	return p
}
func (p *Property) SetMust(value bool) *Property {
	p.Must = value
	return p
}

func (p *Property) SetMin(value int) *Property {
	p.Min = value
	return p
}
func (p *Property) SetMax(value int) *Property {
	p.Max = value
	return p
}
