package google

import (
	"github.com/kubemq-hub/kubemq-source-connectors/config"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
)

func Test_parseOptions(t *testing.T) {
	dat, err := ioutil.ReadFile("./../../../credentials/projectID.txt")
	require.NoError(t, err)
	projectID := string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/subID.txt")
	require.NoError(t, err)
	subID:= string(dat)
	tests := []struct {
		name    string
		cfg     config.Metadata
		want    options
		wantErr bool
	}{
		{
			name: "valid options",
			cfg: config.Metadata{
				Name: "google-pub-sub-source",
				Kind: "",
				Properties: map[string]string{
					"concurrency":            "1",
					"project_id":projectID,
					"sub_id":subID,
				},
			},
			want: options{
				concurrency:         1,
				projectID: projectID,
				subID: subID,
			},
			wantErr: false,
		},
		{
			name: "invalid options - bad concurrency",
			cfg: config.Metadata{
				Name: "google-pub-sub",
				Kind: "",
				Properties: map[string]string{
					"concurrency":            "10000",
					"project_id":projectID,
					"sub_id":subID,
				},
			},
			want:    options{},
			wantErr: true,
		},
		{
			name: "invalid options - bad projectID",
			cfg: config.Metadata{
				Name: "google-pub-sub",
				Kind: "",
				Properties: map[string]string{
					"concurrency":            "10000",
					"sub_id":subID,
				},
			},
			want:    options{},
			wantErr: true,
		},
		{
			name: "invalid options - bad subID ",
			cfg: config.Metadata{
				Name: "google-pub-sub",
				Kind: "",
				Properties: map[string]string{
					"concurrency":            "10000",
					"project_id":projectID,
				},
			},
			want:    options{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseOptions(tt.cfg)
			if tt.wantErr {
				require.Error(t, err)
				
			} else {
				require.NoError(t, err)
				
			}
			
			require.EqualValues(t, got, tt.want)
		})
	}
}


