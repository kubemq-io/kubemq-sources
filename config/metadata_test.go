package config

import (
	"reflect"
	"testing"
)

func TestMetadata_MustParseBool(t *testing.T) {
	type fields struct {
		Name       string
		Kind       string
		Properties map[string]string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid parse",
			fields: fields{
				Name: "",
				Kind: "",
				Properties: map[string]string{
					"bool-key": "true",
				},
			},
			args: args{
				key: "bool-key",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid parse",
			fields: fields{
				Name: "",
				Kind: "",
				Properties: map[string]string{
					"bool-key": "no bool value",
				},
			},
			args: args{
				key: "bool-key",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid parse - no key",
			fields: fields{
				Name: "",
				Kind: "",
				Properties: map[string]string{
					"bool-key": "true",
				},
			},
			args: args{
				key: "not-bool-key",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Metadata{
				Name:       tt.fields.Name,
				Kind:       tt.fields.Kind,
				Properties: tt.fields.Properties,
			}
			got, err := m.MustParseBool(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("MustParseBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MustParseBool() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetadata_MustParseInt(t *testing.T) {
	type fields struct {
		Name       string
		Kind       string
		Properties map[string]string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "valid parse",
			fields: fields{
				Name: "",
				Kind: "",
				Properties: map[string]string{
					"int-key": "1",
				},
			},
			args: args{
				key: "int-key",
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "invalid parse",
			fields: fields{
				Name: "",
				Kind: "",
				Properties: map[string]string{
					"int-key": "bad-value",
				},
			},
			args: args{
				key: "int-key",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "invalid parse - no valid key",
			fields: fields{
				Name: "",
				Kind: "",
				Properties: map[string]string{
					"int-key": "1",
				},
			},
			args: args{
				key: "no-int-key",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Metadata{
				Name:       tt.fields.Name,
				Kind:       tt.fields.Kind,
				Properties: tt.fields.Properties,
			}
			got, err := m.MustParseInt(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("MustParseInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MustParseInt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetadata_MustParseIntWithRange(t *testing.T) {
	type fields struct {
		Name       string
		Kind       string
		Properties map[string]string
	}
	type args struct {
		key string
		min int
		max int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "valid parse",
			fields: fields{
				Name: "",
				Kind: "",
				Properties: map[string]string{
					"int-key": "1",
				},
			},
			args: args{
				key: "int-key",
				min: 1,
				max: 1000,
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "invalid parse",
			fields: fields{
				Name: "",
				Kind: "",
				Properties: map[string]string{
					"int-key": "bad-value",
				},
			},
			args: args{
				key: "int-key",
				min: 1,
				max: 1000,
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "invalid parse - no valid key",
			fields: fields{
				Name: "",
				Kind: "",
				Properties: map[string]string{
					"int-key": "1",
				},
			},
			args: args{
				key: "no-int-key",
				min: 1,
				max: 1000,
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "invalid parse -lower than min",
			fields: fields{
				Name: "",
				Kind: "",
				Properties: map[string]string{
					"int-key": "0",
				},
			},
			args: args{
				key: "int-key",
				min: 1,
				max: 1000,
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "invalid parse -higher than max",
			fields: fields{
				Name: "",
				Kind: "",
				Properties: map[string]string{
					"int-key": "2000",
				},
			},
			args: args{
				key: "int-key",
				min: 1,
				max: 1000,
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Metadata{
				Name:       tt.fields.Name,
				Kind:       tt.fields.Kind,
				Properties: tt.fields.Properties,
			}
			got, err := m.MustParseIntWithRange(tt.args.key, tt.args.min, tt.args.max)
			if (err != nil) != tt.wantErr {
				t.Errorf("MustParseIntWithRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MustParseIntWithRange() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetadata_MustParseString(t *testing.T) {
	type fields struct {
		Name       string
		Kind       string
		Properties map[string]string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "valid parse",
			fields: fields{
				Name: "",
				Kind: "",
				Properties: map[string]string{
					"string-key": "string",
				},
			},
			args: args{
				key: "string-key",
			},
			want:    "string",
			wantErr: false,
		},
		{
			name: "invalid parse",
			fields: fields{
				Name: "",
				Kind: "",
				Properties: map[string]string{
					"string-key": "",
				},
			},
			args: args{
				key: "string-key",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Metadata{
				Name:       tt.fields.Name,
				Kind:       tt.fields.Kind,
				Properties: tt.fields.Properties,
			}
			got, err := m.MustParseString(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("MustParseString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MustParseString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetadata_ParseBool(t *testing.T) {
	type fields struct {
		Name       string
		Kind       string
		Properties map[string]string
	}
	type args struct {
		key          string
		defaultValue bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "valid parse",
			fields: fields{
				Name: "",
				Kind: "",
				Properties: map[string]string{
					"bool-key": "true",
				},
			},
			args: args{
				key:          "bool-key",
				defaultValue: false,
			},
			want: true,
		},
		{
			name: "valid parse with default",
			fields: fields{
				Name: "",
				Kind: "",
				Properties: map[string]string{
					"bool-key": "true",
				},
			},
			args: args{
				key:          "not-bool-key",
				defaultValue: true,
			},
			want: true,
		},
		{
			name: "valid parse with bad input",
			fields: fields{
				Name: "",
				Kind: "",
				Properties: map[string]string{
					"bool-key": "bad-bool-value",
				},
			},
			args: args{
				key:          "bool-key",
				defaultValue: true,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Metadata{
				Name:       tt.fields.Name,
				Kind:       tt.fields.Kind,
				Properties: tt.fields.Properties,
			}
			if got := m.ParseBool(tt.args.key, tt.args.defaultValue); got != tt.want {
				t.Errorf("ParseBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetadata_ParseInt(t *testing.T) {
	type fields struct {
		Name       string
		Kind       string
		Properties map[string]string
	}
	type args struct {
		key          string
		defaultValue int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "valid parse",
			fields: fields{
				Name: "",
				Kind: "",
				Properties: map[string]string{
					"int-key": "1",
				},
			},
			args: args{
				key:          "int-key",
				defaultValue: 100,
			},
			want: 1,
		},
		{
			name: "valid parse with default",
			fields: fields{
				Name: "",
				Kind: "",
				Properties: map[string]string{
					"int-key": "100",
				},
			},
			args: args{
				key:          "not-bool-key",
				defaultValue: 200,
			},
			want: 200,
		},
		{
			name: "valid parse with bad input",
			fields: fields{
				Name: "",
				Kind: "",
				Properties: map[string]string{
					"int-key": "bad-int-value",
				},
			},
			args: args{
				key:          "int-key",
				defaultValue: 300,
			},
			want: 300,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Metadata{
				Name:       tt.fields.Name,
				Kind:       tt.fields.Kind,
				Properties: tt.fields.Properties,
			}
			if got := m.ParseInt(tt.args.key, tt.args.defaultValue); got != tt.want {
				t.Errorf("ParseInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetadata_ParseIntWithRange(t *testing.T) {
	type fields struct {
		Name       string
		Kind       string
		Properties map[string]string
	}
	type args struct {
		key          string
		defaultValue int
		min          int
		max          int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "valid parse",
			fields: fields{
				Name: "",
				Kind: "",
				Properties: map[string]string{
					"int-key": "1",
				},
			},
			args: args{
				key:          "int-key",
				defaultValue: 100,
				min:          1,
				max:          100,
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "invalid parse - lower than min",
			fields: fields{
				Name: "",
				Kind: "",
				Properties: map[string]string{
					"int-key": "0",
				},
			},
			args: args{
				key:          "int-key",
				defaultValue: 100,
				min:          1,
				max:          100,
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "invalid parse - higher than max",
			fields: fields{
				Name: "",
				Kind: "",
				Properties: map[string]string{
					"int-key": "200",
				},
			},
			args: args{
				key:          "int-key",
				defaultValue: 100,
				min:          1,
				max:          100,
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Metadata{
				Name:       tt.fields.Name,
				Kind:       tt.fields.Kind,
				Properties: tt.fields.Properties,
			}
			got, err := m.ParseIntWithRange(tt.args.key, tt.args.defaultValue, tt.args.min, tt.args.max)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseIntWithRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseIntWithRange() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetadata_ParseString(t *testing.T) {
	type fields struct {
		Name       string
		Kind       string
		Properties map[string]string
	}
	type args struct {
		key          string
		defaultValue string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Metadata{
				Name:       tt.fields.Name,
				Kind:       tt.fields.Kind,
				Properties: tt.fields.Properties,
			}
			if got := m.ParseString(tt.args.key, tt.args.defaultValue); got != tt.want {
				t.Errorf("ParseString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetadata_Validate(t *testing.T) {
	type fields struct {
		Name       string
		Kind       string
		Properties map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Metadata{
				Name:       tt.fields.Name,
				Kind:       tt.fields.Kind,
				Properties: tt.fields.Properties,
			}
			if err := m.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMetadata_ParseString1(t *testing.T) {
	type fields struct {
		Name       string
		Kind       string
		Properties map[string]string
	}
	type args struct {
		key          string
		defaultValue string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "valid parse",
			fields: fields{
				Name: "",
				Kind: "",
				Properties: map[string]string{
					"string-key": "string",
				},
			},
			args: args{
				key:          "string-key",
				defaultValue: "default",
			},
			want: "string",
		},
		{
			name: "valid parse with default",
			fields: fields{
				Name: "",
				Kind: "",
				Properties: map[string]string{
					"string-key": "string",
				},
			},
			args: args{
				key:          "not-string-key",
				defaultValue: "default",
			},
			want: "default",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Metadata{
				Name:       tt.fields.Name,
				Kind:       tt.fields.Kind,
				Properties: tt.fields.Properties,
			}
			if got := m.ParseString(tt.args.key, tt.args.defaultValue); got != tt.want {
				t.Errorf("ParseString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetadata_MustParseMap(t *testing.T) {
	type fields struct {
		Name       string
		Kind       string
		Properties map[string]string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "valid parse",
			fields: fields{
				Name: "",
				Kind: "",
				Properties: map[string]string{
					"map-key": `{ "key":"value"}`,
				},
			},
			args: args{
				key: "map-key",
			},
			want:    map[string]string{"key": "value"},
			wantErr: false,
		},
		{
			name: "invalid parse",
			fields: fields{
				Name: "",
				Kind: "",
				Properties: map[string]string{
					"map-key": `{ "key":value}`,
				},
			},
			args: args{
				key: "map-key",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid parse - empty ",
			fields: fields{
				Name: "",
				Kind: "",
				Properties: map[string]string{
					"map-key": ``,
				},
			},
			args: args{
				key: "map-key",
			},
			want:    map[string]string{},
			wantErr: false,
		},
		{
			name: "valid parse - no key ",
			fields: fields{
				Name: "",
				Kind: "",
				Properties: map[string]string{
					"map-key": ``,
				},
			},
			args: args{
				key: "no-map-key",
			},
			want:    map[string]string{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Metadata{
				Name:       tt.fields.Name,
				Kind:       tt.fields.Kind,
				Properties: tt.fields.Properties,
			}
			got, err := m.MustParseJsonMap(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("MustParseJsonMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MustParseJsonMap() got = %v, want %v", got, tt.want)
			}
		})
	}
}
