package sqs

import (
	"github.com/kubemq-hub/kubemq-source-connectors/config"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_parseOptions(t *testing.T) {
	aswKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecret := os.Getenv("AWS_SECRET_ACCESS_KEY")
	sqsQueue := os.Getenv("SQS_QUEUE_NAME")
	tests := []struct {
		name    string
		cfg     config.Metadata
		want    options
		wantErr bool
	}{
		{
			name: "valid options",
			cfg: config.Metadata{
				Name: "sqs-queue",
				Kind: "",
				Properties: map[string]string{
					"sqs_key":                aswKey,
					"sqs_secret_key":         awsSecret,
					"queue":                  sqsQueue,
					"region":                 "us-west-2",
					"visibility":             "0",
					"max_number_of_messages": "0",
					"wait_time_seconds":      "0",
					"concurrency":            "1",
				},
			},
			want: options{
				sqsKey:              aswKey,
				sqsSecretKey:        awsSecret,
				region:              "us-west-2",
				visibility:          0,
				queue:               sqsQueue,
				maxNumberOfMessages: 0,
				concurrency:         1,
				waitTimeSeconds:     0,
			},
			wantErr: false,
		},
		{
			name: "invalid options - bad concurrency",
			cfg: config.Metadata{
				Name: "sqs-queue",
				Kind: "",
				Properties: map[string]string{
					"sqs_key":                aswKey,
					"sqs_secret_key":         awsSecret,
					"queue":                  sqsQueue,
					"region":                 "us-west-2",
					"visibility":             "0",
					"max_number_of_messages": "0",
					"wait_time_seconds":      "0",
					"concurrency":            "10000",
				},
			},
			want:    options{},
			wantErr: true,
		},
		{
			name: "invalid options - bad channel",
			cfg: config.Metadata{
				Name: "kubemq-rpc",
				Kind: "",
				Properties: map[string]string{
					"host":    "localhost",
					"port":    "50000",
					"channel": "",
				},
			},
			want:    options{},
			wantErr: true,
		},
		{
			name: "invalid options - bad sqs secret",
			cfg: config.Metadata{
				Name: "sqs-queue",
				Kind: "",
				Properties: map[string]string{
					"sqs_key":                aswKey,
					"queue":                  sqsQueue,
					"region":                 "us-west-2",
					"visibility":             "0",
					"max_number_of_messages": "0",
					"wait_time_seconds":      "0",
					"concurrency":            "10000",
				},
			},
			want:    options{},
			wantErr: true,
		},{
			name: "invalid options - bad sqs region",
			cfg: config.Metadata{
				Name: "sqs-queue",
				Kind: "",
				Properties: map[string]string{
					"sqs_key":                aswKey,
					"sqs_secret_key":         awsSecret,
					"queue":                  sqsQueue,
					"visibility":             "0",
					"max_number_of_messages": "0",
					"wait_time_seconds":      "0",
					"concurrency":            "10000",
				},
			},
			want:    options{},
			wantErr: true,
		},
		{
			name: "invalid options - bad queue",
			cfg: config.Metadata{
				Name: "kubemq-rpc",
				Kind: "",
				Properties: map[string]string{
					"sqs_key":                aswKey,
					"sqs_secret_key":         awsSecret,
					"region":                 "us-west-2",
					"visibility":             "0",
					"max_number_of_messages": "0",
					"wait_time_seconds":      "0",
					"concurrency":            "1",
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


