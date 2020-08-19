package sqs
//
//import (
//	"context"
//	"github.com/kubemq-hub/kubemq-sources/targets"
//	"github.com/kubemq-hub/kubemq-sources/targets/null"
//	"github.com/kubemq-hub/kubemq-sources/types"
//	"github.com/stretchr/testify/require"
//	"os"
//
//	"github.com/kubemq-hub/kubemq-sources/config"
//
//	"testing"
//	"time"
//)
//
//func TestClient_Init(t *testing.T) {
//
//	aswKey := os.Getenv("AWS_ACCESS_KEY_ID")
//	awsSecret := os.Getenv("AWS_SECRET_ACCESS_KEY")
//	sqsQueue := os.Getenv("SQS_QUEUE_NAME")
//	deadLetter := os.Getenv("DEAD_LETTER")
//	tests := []struct {
//		name    string
//		cfg     config.Metadata
//		wantErr bool
//	}{
//		{
//			name: "init",
//			cfg: config.Metadata{
//				Name: "sqs-queue",
//				Kind: "",
//				Properties: map[string]string{
//					"sqs_key":                aswKey,
//					"sqs_secret_key":         awsSecret,
//					"queue":                  sqsQueue,
//					"region":                 "us-west-2",
//					"visibility":             "0",
//					"dead_letter":             deadLetter,
//					"max_receive":             "2",
//					"max_number_of_messages": "0",
//					"wait_time_seconds":      "0",
//					"concurrency":            "1",
//				},
//			},
//			wantErr: false,
//		},
//		{
//			name: "init - error",
//			cfg: config.Metadata{
//				Name: "sqs-queue",
//				Kind: "",
//				Properties: map[string]string{
//					"region":                 "us-west-2",
//					"visibility":             "0",
//					"max_number_of_messages": "0",
//					"wait_time_seconds":      "0",
//					"concurrency":            "1",
//				},
//			},
//			wantErr: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//			defer cancel()
//			c := New()
//			err := c.Init(ctx, tt.cfg)
//			if tt.wantErr {
//				require.Error(t, err)
//				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
//				return
//			}
//			require.NoError(t, err)
//		})
//	}
//}
//
//func TestClient_Start(t *testing.T) {
//	aswKey := os.Getenv("AWS_ACCESS_KEY_ID")
//	awsSecret := os.Getenv("AWS_SECRET_ACCESS_KEY")
//	sqsQueue := os.Getenv("SQS_QUEUE_NAME")
//	tests := []struct {
//		name    string
//		target  -sources.Target
//		cfg     config.Metadata
//		wantErr bool
//	}{
//		{
//			name: "start",
//			target: &null.Client{
//				Delay:         0,
//				DoError:       nil,
//				ResponseError: nil,
//			},
//			cfg: config.Metadata{
//				Name: "sqs-queue",
//				Kind: "",
//				Properties: map[string]string{
//					"sqs_key":                aswKey,
//					"sqs_secret_key":         awsSecret,
//					"queue":                  sqsQueue,
//					"region":                 "us-west-2",
//					"visibility":             "0",
//					"max_number_of_messages": "2",
//					"wait_time_seconds":      "0",
//					"concurrency":            "1",
//				},
//			},
//			wantErr: false,
//		},
//		{
//			name:   "start - bad target",
//			target: nil,
//			cfg: config.Metadata{
//				Name: "sqs-queue",
//				Kind: "",
//				Properties: map[string]string{
//					"sqs_key":                aswKey,
//					"sqs_secret_key":         awsSecret,
//					"queue":                  "sqsQueue",
//					"region":                 "us-west-2",
//					"visibility":             "0",
//					"max_number_of_messages": "0",
//					"wait_time_seconds":      "0",
//					"concurrency":            "1",
//				},
//			},
//			wantErr: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//			defer cancel()
//			c := New()
//			err := c.Init(ctx, tt.cfg)
//			require.NoError(t, err)
//			err = c.Start(ctx, tt.target)
//			if tt.wantErr {
//				require.Error(t, err)
//				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
//				return
//			}
//			//For debugging
//			time.Sleep(10 * time.Millisecond)
//			require.NoError(t, err)
//		})
//	}
//}
//
//func TestClient_SetQueueAttributes(t *testing.T) {
//	aswKey := os.Getenv("AWS_ACCESS_KEY_ID")
//	awsSecret := os.Getenv("AWS_SECRET_ACCESS_KEY")
//	sqsQueue := os.Getenv("SQS_QUEUE_NAME")
//	deadLetter:= os.Getenv("DEAD_LETTER_QUEUE")
//	tests := []struct {
//		name    string
//		cfg     config.Metadata
//		queueURL string
//		want    *types.Response
//		wantErr bool
//	}{
//		{
//			name: "valid set queue attribute",
//			cfg: config.Metadata{
//				Name: "target.sqs",
//				Kind: "target.sqs",
//				Properties: map[string]string{
//					"sqs_key":                     aswKey,
//					"sqs_secret_key":              awsSecret,
//					"region":                      "us-west-2",
//					"max_receive":                 "10",
//					"dead_letter":                 deadLetter,
//					"max_retries":                 "0",
//					"max_retries_backoff_seconds": "0",
//					"retries":                     "0",
//				},
//			},
//			queueURL: sqsQueue,
//			wantErr: false,
//		},{
//			name: "in-valid set queue attribute",
//			cfg: config.Metadata{
//				Name: "target.sqs",
//				Kind: "target.sqs",
//				Properties: map[string]string{
//					"sqs_key":                     aswKey,
//					"sqs_secret_key":              awsSecret,
//					"region":                      "us-west-2",
//					"dead_letter":                 deadLetter,
//					"max_retries":                 "0",
//					"max_retries_backoff_seconds": "0",
//					"retries":                     "0",
//				},
//			},
//			queueURL: sqsQueue,
//			wantErr: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ctx, cancel := context.WithCancel(context.Background())
//			defer cancel()
//			c := New()
//			err := c.Init(ctx, tt.cfg)
//			require.NoError(t, err)
//			err = c.SetQueueAttributes(ctx, sqsQueue)
//			if tt.wantErr {
//				require.Error(t, err)
//				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
//				return
//			}
//			require.NoError(t, err)
//		})
//	}
//}
