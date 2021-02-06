package filesystem

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/kubemq-hub/kubemq-sources/middleware"
	"github.com/kubemq-hub/kubemq-sources/pkg/logger"

	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	pollInterval = 5 * time.Second
)

type Client struct {
	opts       options
	waiting    sync.Map
	inProgress sync.Map
	completed  sync.Map
	sendCh     chan *SourceFile
	logger     *logger.Logger
	ctx        context.Context
	cancelFunc context.CancelFunc
	absRoot    string
}

func New() *Client {
	return &Client{
		waiting:    sync.Map{},
		inProgress: sync.Map{},
		completed:  sync.Map{},
	}
}
func (c *Client) Connector() *common.Connector {
	return Connector()
}
func (c *Client) Init(ctx context.Context, cfg config.Spec) error {
	c.logger = logger.NewLogger(cfg.Name)
	var err error
	c.opts, err = parseOptions(cfg)
	if err != nil {
		return err
	}
	c.absRoot, _ = filepath.Abs(filepath.Clean(c.opts.root))
	return nil
}

func (c *Client) Start(ctx context.Context, target middleware.Middleware) error {
	if _, err := os.Stat(c.absRoot); os.IsNotExist(err) {
		return fmt.Errorf("root %s path is not exist", c.absRoot)
	}
	c.ctx, c.cancelFunc = context.WithCancel(ctx)
	c.sendCh = make(chan *SourceFile)
	go c.scan(c.ctx)
	for i := 0; i < c.opts.concurrency; i++ {
		go c.senderFunc(c.ctx, target)
	}
	go c.send(c.ctx)
	return nil
}
func (c *Client) Stop() error {
	c.cancelFunc()
	c.waiting = sync.Map{}
	return nil
}
func (c *Client) inPipe(file *SourceFile) bool {
	if _, ok := c.waiting.Load(file.FullPath()); ok {
		return true
	}
	if _, ok := c.inProgress.Load(file.FullPath()); ok {
		return true
	}
	if _, ok := c.completed.Load(file.FullPath()); ok {
		c.logger.Infof("file %s already sent and will be deleted", file.FullPath())
		if err := file.Delete(); err != nil {
			c.logger.Errorf("error during delete a file %s,%s, will try again", file.FullPath(), err.Error())
		}
		return true
	}
	return false
}
func (c *Client) walk() error {
	var list []*SourceFile
	err := filepath.Walk(c.absRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			list = append(list, NewSourceFile(info, path, c.absRoot))
		}
		return nil
	})

	if err != nil {
		return err
	}
	added := 0
	for _, file := range list {
		if !c.inPipe(file) {
			c.waiting.Store(file.FullPath(), file)
			added++
		}
	}
	if added > 0 {
		c.logger.Infof("%d new files added to sending waiting list", added)
	}
	return nil
}

func (c *Client) senderFunc(ctx context.Context, sender middleware.Middleware) {
	for {
		select {
		case file := <-c.sendCh:
			c.inProgress.Store(file.FullPath(), file)
			c.waiting.Delete(file.FullPath())
			req, err := file.Request(c.opts.bucketType, c.opts.bucketName)
			if err != nil {
				c.logger.Errorf("error during creating file requests %s, %s", file.FullPath(), err.Error())
				c.waiting.Store(file.FullPath(), file)
				c.inProgress.Delete(file.FullPath())
				continue
			}
			c.logger.Infof("sending file %s started", file.FileName())
			resp, err := sender.Do(ctx, req)
			if err != nil {
				c.logger.Errorf("error during sending file %s, %s", file.FileName(), err.Error())
				c.waiting.Store(file.FullPath(), file)
				c.inProgress.Delete(file.FullPath())
				continue
			}
			if resp.IsError {
				c.logger.Errorf("error on sending file %s response, %s", file.FileName(), resp.Error)
				c.waiting.Store(file.FullPath(), file)
				c.inProgress.Delete(file.FullPath())
				continue
			}
			if err := file.Delete(); err != nil {
				c.logger.Errorf("error during delete a file %s, %s,file will be resend", file.FileName(), err.Error())
				c.waiting.Store(file.FullPath(), file)
			} else {
				c.completed.Store(file.FullPath(), file)
			}
			c.inProgress.Delete(file.FullPath())
			c.logger.Infof("sending file %s completed", file.FileName())
		case <-ctx.Done():
			return
		}
	}
}
func (c *Client) scan(ctx context.Context) {
	for {
		select {
		case <-time.After(pollInterval):
			err := c.walk()
			if err != nil {
				c.logger.Errorf("error during scan files, %s", err.Error())
			}
		case <-ctx.Done():
			return
		}
	}
}
func (c *Client) send(ctx context.Context) {
	for {
		select {
		case <-time.After(time.Second):
			var list []*SourceFile
			c.waiting.Range(func(key, value interface{}) bool {
				list = append(list, value.(*SourceFile))
				return true
			})
			for _, file := range list {
				if _, ok := c.inProgress.Load(file.FullPath()); ok {
					continue
				}
				select {
				case c.sendCh <- file:
				case <-ctx.Done():
					return
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
