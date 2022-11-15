package http

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemq-sources/config"
	targetMiddleware "github.com/kubemq-io/kubemq-sources/middleware"
	"github.com/kubemq-io/kubemq-sources/pkg/logger"
	"github.com/kubemq-io/kubemq-sources/types"
)

type Client struct {
	log     *logger.Logger
	opts    options
	target  targetMiddleware.Middleware
	Methods []string
	Path    string
}

func (c *Client) Connector() *common.Connector {
	return Connector()
}

func New() *Client {
	return &Client{}
}

func (c *Client) Init(ctx context.Context, cfg config.Spec, log *logger.Logger) error {
	c.log = log
	if c.log == nil {
		c.log = logger.NewLogger(cfg.Kind)
	}
	var err error
	c.opts, err = parseOptions(cfg)
	if err != nil {
		return err
	}
	c.Methods = c.opts.methods
	c.Path = c.opts.path
	return nil
}

func (c *Client) Start(ctx context.Context, target targetMiddleware.Middleware) error {
	c.target = target
	return nil
}

func (c *Client) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	req, err := c.parseRequest(r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	resp, err := c.target.Do(ctx, req)
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

func (c *Client) Stop() error {
	return nil
}

func (c *Client) parseRequest(httpRequest *http.Request) (*types.Request, error) {
	mdBuff, err := httputil.DumpRequest(httpRequest, false)
	if err != nil {
		return nil, err
	}
	req := types.NewRequest().SetMetadata(string(mdBuff))
	if c.opts.dynamicMapping {
		req.SetChannel(strings.Replace(httpRequest.URL.Path, "/", ".", -1))
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
