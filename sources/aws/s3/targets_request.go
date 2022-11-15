package s3

import (
	"encoding/json"
	"fmt"
)

type TargetsRequest struct {
	Metadata TargetsMetadata `json:"metadata,omitempty"`
	Data     []byte          `json:"data,omitempty"`
}

func NewTargetsRequest() *TargetsRequest {
	return &TargetsRequest{
		Metadata: NewTargetMetadata(),
		Data:     nil,
	}
}

func (r *TargetsRequest) SetMetadata(value TargetsMetadata) *TargetsRequest {
	r.Metadata = value
	return r
}

func (r *TargetsRequest) SetMetadataKeyValue(key, value string) *TargetsRequest {
	r.Metadata.Set(key, value)
	return r
}

func (r *TargetsRequest) SetData(value []byte) *TargetsRequest {
	r.Data = value
	return r
}

func (r *TargetsRequest) Size() float64 {
	return float64(len(r.Data))
}

func ParseRequest(body []byte) (*TargetsRequest, error) {
	if body == nil {
		return nil, fmt.Errorf("empty request")
	}
	req := &TargetsRequest{}
	err := json.Unmarshal(body, req)
	if err != nil {
		return NewTargetsRequest().SetData(body), err
	}
	return req, nil
}

func (r *TargetsRequest) MarshalBinary() []byte {
	data, _ := json.Marshal(r)
	return data
}
