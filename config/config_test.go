package config

import (
	"os"
	"reflect"
	"testing"
)

func TestConfig_Validate(t *testing.T) {

	tests := []struct {
		name     string
		Bindings []BindingConfig
		wantErr  bool
	}{
		{
			name: "valid config",
			Bindings: []BindingConfig{
				{
					Name: "binding-1",
					Source: Spec{
						Name:       "source-1",
						Kind:       "source-1",
						Properties: nil,
					},

					Target: Spec{

						Name:       "target-1",
						Kind:       "target-1",
						Properties: nil,
					},
				},
			},
			wantErr: false,
		},
		{
			name:     "invalid config - no bindings",
			Bindings: []BindingConfig{},
			wantErr:  true,
		},
		{
			name: "invalid config - binding no name",
			Bindings: []BindingConfig{
				{
					Source: Spec{
						Name:       "source-1",
						Kind:       "source-1",
						Properties: nil,
					},

					Target: Spec{

						Name:       "target-1",
						Kind:       "target-1",
						Properties: nil,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid config - invalid source",
			Bindings: []BindingConfig{
				{
					Name: "binding-1",
					Source: Spec{
						Kind:       "source-1",
						Properties: nil,
					},

					Target: Spec{

						Name:       "target-1",
						Kind:       "target-1",
						Properties: nil,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid config - bad target",
			Bindings: []BindingConfig{
				{
					Name: "binding-1",
					Source: Spec{
						Name:       "source-1",
						Kind:       "source-1",
						Properties: nil,
					},

					Target: Spec{
						Kind:       "target-1",
						Properties: nil,
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Bindings: tt.Bindings,
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
bindings:
- name: binding-1
  source: 
    kind: source-1
    name: source-1
    properties: null
  target: 
    kind: target-1
    name: target-1
    properties: null	
`,
			want: &Config{
				Bindings: []BindingConfig{
					{
						Name: "binding-1",
						Source: Spec{
							Name:       "source-1",
							Kind:       "source-1",
							Properties: nil,
						},

						Target: Spec{

							Name:       "target-1",
							Kind:       "target-1",
							Properties: nil,
						},
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
