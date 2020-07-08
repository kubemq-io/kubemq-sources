package config

import (
	"os"
	"reflect"
	"testing"
)

func TestConfig_Validate(t *testing.T) {
	type fields struct {
		Sources  []Metadata
		Targets  []Metadata
		Bindings []Binding
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid metadata",
			fields: fields{
				Sources: []Metadata{
					{
						Name:       "source-1",
						Kind:       "source-1",
						Properties: nil,
					},
				},
				Targets: []Metadata{
					{
						Name:       "target-1",
						Kind:       "target-1",
						Properties: nil,
					},
				},
				Bindings: []Binding{
					{
						Source: "source-1",
						Target: "target-1",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid metadata - sources error",
			fields: fields{
				Sources: []Metadata{
					{
						Name:       "",
						Kind:       "source-1",
						Properties: nil,
					},
				},
				Targets: []Metadata{
					{
						Name:       "target-1",
						Kind:       "target-1",
						Properties: nil,
					},
				},
				Bindings: []Binding{
					{
						Source: "source-1",
						Target: "target-1",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid metadata - target error",
			fields: fields{
				Sources: []Metadata{
					{
						Name:       "source-1",
						Kind:       "source-1",
						Properties: nil,
					},
				},
				Targets: []Metadata{
					{
						Name:       "",
						Kind:       "target-1",
						Properties: nil,
					},
				},
				Bindings: []Binding{
					{
						Source: "source-1",
						Target: "target-1",
					},
				},
			},
			wantErr: true,
		},

		{
			name: "invalid metadata - no sources error",
			fields: fields{
				Sources: nil,
				Targets: []Metadata{
					{
						Name:       "target-1",
						Kind:       "target-1",
						Properties: nil,
					},
				},
				Bindings: []Binding{
					{
						Source: "source-1",
						Target: "target-1",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid metadata - no targets error",
			fields: fields{
				Sources: []Metadata{
					{
						Name:       "source-1",
						Kind:       "source-1",
						Properties: nil,
					},
				},
				Targets: nil,
				Bindings: []Binding{
					{
						Source: "source-1",
						Target: "target-1",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid metadata - no binding error",
			fields: fields{
				Sources: []Metadata{
					{
						Name:       "source-1",
						Kind:       "source-1",
						Properties: nil,
					},
				},
				Targets: []Metadata{
					{
						Name:       "target-1",
						Kind:       "target-1",
						Properties: nil,
					},
				},
				Bindings: nil,
			},
			wantErr: true,
		},
		{
			name: "invalid metadata - binding error entry 1",
			fields: fields{
				Sources: []Metadata{
					{
						Name:       "source-1",
						Kind:       "source-1",
						Properties: nil,
					},
				},
				Targets: []Metadata{
					{
						Name:       "target-1",
						Kind:       "target-1",
						Properties: nil,
					},
				},
				Bindings: []Binding{
					{
						Source: "",
						Target: "target-1",
					},
				},
			},
			wantErr: true,
		}, {
			name: "invalid metadata - binding error entry 2",
			fields: fields{
				Sources: []Metadata{
					{
						Name:       "source-1",
						Kind:       "source-1",
						Properties: nil,
					},
				},
				Targets: []Metadata{
					{
						Name:       "target-1",
						Kind:       "target-1",
						Properties: nil,
					},
				},
				Bindings: []Binding{
					{
						Source: "source-1",
						Target: "",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Sources:  tt.fields.Sources,
				Targets:  tt.fields.Targets,
				Bindings: tt.fields.Bindings,
			}
			if err := c.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			//y, _ := yaml.Marshal(c)
			//
			//ioutil.WriteFile("config.yaml", y, 0644)
		})
	}
}

func TestLoad_Env_Yaml(t *testing.T) {

	tests := []struct {
		name      string
		cfgString string
		want      *Config
		wantErr   bool
	}{
		{
			name: "load from env",
			cfgString: `
Sources:
- Kind: source-1
  Name: source-1
  Properties: null
Targets:
- Kind: target-1
  Name: target-1
  Properties: null
Bindings:
- Source: source-1
  Target: target-1	
`,
			want: &Config{
				Sources: []Metadata{
					{
						Name:       "source-1",
						Kind:       "source-1",
						Properties: nil,
					},
				},
				Targets: []Metadata{
					{
						Name:       "target-1",
						Kind:       "target-1",
						Properties: nil,
					},
				},
				Bindings: []Binding{
					{
						Source: "source-1",
						Target: "target-1",
					},
				},
			},
			wantErr: false,
		},
		{
			name:      "invalid file format",
			cfgString: "some-invalid-format",
			want:      nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.Setenv("CONFIG", tt.cfgString)
			defer os.RemoveAll("./config.yaml")
			got, err := Load()
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() got = %v, want %v", got, tt.want)
			}
		})
	}
}
