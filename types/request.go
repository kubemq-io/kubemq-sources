package types

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/sqs"
	jsoniter "github.com/json-iterator/go"
	"github.com/kubemq-io/kubemq-go"
	"cloud.google.com/go/pubsub"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Request struct {
	Metadata Metadata `json:"metadata"`
	Data     []byte   `json:"data"`
}

func NewRequest() *Request {
	return &Request{
		Metadata: NewMetadata(),
		Data:     nil,
	}
}

func (r *Request) SetMetadata(value Metadata) *Request {
	r.Metadata = value
	return r
}
func (r *Request) SetMetadataKeyValue(key, value string) *Request {
	r.Metadata.Set(key, value)
	return r
}

func (r *Request) SetData(value []byte) *Request {
	r.Data = value
	return r
}

func ParseRequest(body []byte) (*Request, error) {
	if body == nil {
		return nil, fmt.Errorf("empty request")
	}
	req := &Request{}
	err := json.Unmarshal(body, req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (r *Request) MarshalBinary() []byte {
	data, _ := json.Marshal(r)
	return data
}
func (r *Request) ToEvent() *kubemq.Event {
	return kubemq.NewEvent().
		SetMetadata(r.Metadata.String()).
		SetBody(r.Data)
}
func (r *Request) ToEventStore() *kubemq.EventStore {
	return kubemq.NewEventStore().
		SetMetadata(r.Metadata.String()).
		SetBody(r.Data)
}
func (r *Request) ToCommand() *kubemq.Command {
	return kubemq.NewCommand().
		SetMetadata(r.Metadata.String()).
		SetBody(r.Data)
}
func (r *Request) ToQuery() *kubemq.Query {
	return kubemq.NewQuery().
		SetMetadata(r.Metadata.String()).
		SetBody(r.Data)
}
func (r *Request) ToQueueMessage() *kubemq.QueueMessage {
	return kubemq.NewQueueMessage().
		SetMetadata(r.Metadata.String()).
		SetBody(r.Data)
}

func parseRequest(meta string, body []byte) (*Request, error) {
	if meta == "{}" && body == nil {
		return nil, fmt.Errorf("error request is empty")
	}
	req := NewRequest()
	parsedMeta, err := UnmarshallMetadata(meta)
	if err != nil {
		return nil, fmt.Errorf("error parsing request metadata, %w", err)
	}
	return req.
			SetMetadata(parsedMeta).
			SetData(body),
		nil
}
func ParseRequestFromEvent(event *kubemq.Event) (*Request, error) {
	return parseRequest(event.Metadata, event.Body)
}

func ParseRequestFromEventStore(event *kubemq.EventStore) (*Request, error) {
	return parseRequest(event.Metadata, event.Body)
}
func ParseRequestFromEventStoreReceive(event *kubemq.EventStoreReceive) (*Request, error) {
	return parseRequest(event.Metadata, event.Body)
}

func ParseRequestFromCommand(cmd *kubemq.Command) (*Request, error) {
	return parseRequest(cmd.Metadata, cmd.Body)
}
func ParseRequestFromCommandReceive(cmd *kubemq.CommandReceive) (*Request, error) {
	return parseRequest(cmd.Metadata, cmd.Body)
}

func ParseRequestFromQuery(query *kubemq.Query) (*Request, error) {
	return parseRequest(query.Metadata, query.Body)
}

func ParseRequestFromQueryReceive(query *kubemq.QueryReceive) (*Request, error) {
	return parseRequest(query.Metadata, query.Body)
}

func ParseRequestFromQueueMessage(msg *kubemq.QueueMessage) (*Request, error) {
	return parseRequest(msg.Metadata, msg.Body)
}

func ParseRequestFromSQSMessage(msg *sqs.Message) (*Request, error) {
	m := make(map[string]string)
	for k, v := range msg.MessageAttributes {
		m[k] = *v.StringValue
	}
	jsonString, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	req := NewRequest()
	req.SetMetadataKeyValue("tags", string(jsonString))
	return req.
			SetData([]byte(*msg.Body)),
		nil
}

func ParseRequestFromGooglePubSubMessage(msg *pubsub.Message) (*Request, error) {
	m := make(map[string]string)
	for k, v := range msg.Attributes {
		m[k] = v
	}
	jsonString, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	req := NewRequest()
	req.SetMetadataKeyValue("tags", string(jsonString))
	return req.
			SetData(msg.Data),
		nil
}
