package http

import (
	"context"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/kubemq-sources/config"
	targetMiddleware "github.com/kubemq-hub/kubemq-sources/middleware"
	"github.com/kubemq-hub/kubemq-sources/types"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
)

type Handler struct {
	opts    options
	target  targetMiddleware.Middleware
	Methods []string
	Path    string
}

func (h *Handler) Connector() *common.Connector {
	return Connector()
}

func New() *Handler {
	return &Handler{}

}

func (h *Handler) Init(ctx context.Context, cfg config.Spec) error {
	var err error
	h.opts, err = parseOptions(cfg)
	if err != nil {
		return err
	}
	h.Methods = h.opts.methods
	h.Path = h.opts.path
	return nil
}

func (h *Handler) Start(ctx context.Context, target targetMiddleware.Middleware) error {
	h.target = target
	return nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	req, err := h.parseRequest(r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	resp, err := h.target.Do(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if resp.IsError {
		http.Error(w, resp.Error, 500)
		return
	}
	w.WriteHeader(200)
	_, _ = w.Write(resp.Data)

}
func (h *Handler) Stop() error {
	return nil
}

func (h *Handler) parseRequest(httpRequest *http.Request) (*types.Request, error) {
	mdBuff, err := httputil.DumpRequest(httpRequest, false)
	if err != nil {
		return nil, err
	}
	req := types.NewRequest().SetMetadata(string(mdBuff))
	if h.opts.dynamicMapping {
		req.SetChannel(strings.TrimPrefix(httpRequest.RequestURI, "/"))
	}
	if httpRequest.Body == nil {
		return req, nil
	}

	body, err := ioutil.ReadAll(httpRequest.Body)
	if err != nil {
		return nil, err
	}
	req.SetData(body)
	return req, nil
}
