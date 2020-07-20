package google

import (
	"context"
	"github.com/kubemq-hub/kubemq-sources/targets"
	"github.com/kubemq-hub/kubemq-sources/targets/null"
	"github.com/stretchr/testify/require"
	"io/ioutil"

	"github.com/kubemq-hub/kubemq-sources/config"

	"testing"
	"time"
)

func TestClient_Init(t *testing.T) {
	dat, err := ioutil.ReadFile("./../../../credentials/projectID.txt")
	require.NoError(t, err)
	projectID := string(dat)
	require.NoError(t, err)
	dat, err = ioutil.ReadFile("./../../../credentials/subID.txt")
	require.NoError(t, err)
	subID:= string(dat)
	tests := []struct {
		name    string
		cfg     config.Metadata
		wantErr bool
	}{
		{
			name: "init",
			cfg: config.Metadata{
				Name: "google-pub-sub-source",
				Kind: "",
				Properties: map[string]string{
					"max_number_of_messages": "0",
					"concurrency":            "1",
					"project_id":projectID,
					"max_wait_time":"1000",
					"sub_id":subID,
				},
			},
			wantErr: false,
		},
		{
			name: "init-missing-project-id",
			cfg: config.Metadata{
				Name: "google-pub-sub-source",
				Kind: "",
				Properties: map[string]string{
					"max_number_of_messages": "0",
					"concurrency":            "1",
					"sub_id":subID,
				},
			},
			wantErr: true,
		},{
			name: "init-missing-project-sub_id",
			cfg: config.Metadata{
				Name: "google-pub-sub-source",
				Kind: "",
				Properties: map[string]string{
					"max_number_of_messages": "0",
					"concurrency":            "1",
					"project_id":projectID,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
			defer cancel()
			c := New()

			err := c.Init(ctx, tt.cfg)
			if tt.wantErr {
				require.Error(t, err)
				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestClient_Start(t *testing.T) {
	dat, err := ioutil.ReadFile("./../../../credentials/projectID.txt")
	require.NoError(t, err)
	projectID := string(dat)
	require.NoError(t, err)
	dat, err = ioutil.ReadFile("./../../../credentials/subID.txt")
	require.NoError(t, err)
	subID:= string(dat)
	tests := []struct {
		name    string
		target  targets.Target
		cfg     config.Metadata
		wantErr bool
	}{
		{
			name: "init",
			target: &null.Client{
				Delay:         0,
				DoError:       nil,
				ResponseError: nil,
			},
			cfg: config.Metadata{
				Name: "google-pub-sub-source",
				Kind: "",
				Properties: map[string]string{
					"concurrency":            "1",
					"project_id":projectID,
					"sub_id":subID,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
			defer cancel()
			c := New()
			err = c.Init(ctx, tt.cfg)
			require.NoError(t, err)
			err := c.Start(ctx, tt.target)
			if tt.wantErr {
				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
				require.Error(t, err)
				return
			}
			//For debugging
			time.Sleep(100 * time.Second)
			require.NoError(t, err)
		})
	}
}

