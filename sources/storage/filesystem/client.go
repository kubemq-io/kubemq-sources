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

type Client struct {
	opts           options
	waiting        sync.Map
	inProgress     sync.Map
	completed      sync.Map
	sendCh         chan *SourceFile
	log            *logger.Logger
	ctx            context.Context
	cancelFunc     context.CancelFunc
	absRootFolders map[string]string
}

func New() *Client {
	return &Client{
		waiting:        sync.Map{},
		inProgress:     sync.Map{},
		completed:      sync.Map{},
		absRootFolders: map[string]string{},
	}
}
func (c *Client) Connector() *common.Connector {
	return Connector()
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
	for _, folder := range c.opts.folders {
		absPath, _ := filepath.Abs(filepath.Clean(folder))
		c.absRootFolders[absPath] = absPath
	}
	return nil
}

func (c *Client) Start(ctx context.Context, target middleware.Middleware) error {
	for _, folder := range c.absRootFolders {
		if _, err := os.Stat(folder); os.IsNotExist(err) {
			return fmt.Errorf("folder %s path does not exist", folder)
		}
	}
	if c.opts.backupFolder != "" {
		if _, err := os.Stat(c.opts.backupFolder); os.IsNotExist(err) {
			err = os.MkdirAll(c.opts.backupFolder, 0660)
			if err != nil {
				return err
			}
		}
		for _, source := range c.absRootFolders {
			if source == c.opts.backupFolder {
				return fmt.Errorf("move folder %s path cannot be source folder match (recursive)", c.opts.backupFolder)
			}
		}
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
	if val, ok := c.completed.Load(file.FullPath()); ok {
		current := val.(*SourceFile)
		if current.Hash() == file.Hash() {
			c.log.Infof("file %s already sent and will be deleted", file.FullPath())
			if err := file.Delete(); err != nil {
				c.log.Errorf("error during delete a file %s,%s, will try again", file.FullPath(), err.Error())
			}
			return true
		} else {
			c.log.Infof("file %s already sent but a new content has been detected, resending", file.FullPath())
		}
	}
	return false
}
func (c *Client) walk(folder string) error {
	var list []*SourceFile
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			list = append(list, NewSourceFile(info, path, folder, c.opts.backupFolder))
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
		c.log.Debugf("%d new files added to sending waiting list", added)
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
				c.log.Errorf("error during creating file requests %s, %s", file.FullPath(), err.Error())
				c.waiting.Store(file.FullPath(), file)
				c.inProgress.Delete(file.FullPath())
				continue
			}
			c.log.Infof("sending %s started ", file.Metadata())
			resp, err := sender.Do(ctx, req)
			if err != nil {
				c.log.Errorf("error during sending file %s, %s", file.FileName(), err.Error())
				c.waiting.Store(file.FullPath(), file)
				c.inProgress.Delete(file.FullPath())
				continue
			}
			if resp.IsError {
				c.log.Errorf("error on sending file %s response, %s", file.FileName(), resp.Error)
				c.waiting.Store(file.FullPath(), file)
				c.inProgress.Delete(file.FullPath())
				continue
			}
			if err := file.Do(); err != nil {
				c.log.Errorf("error during delete/moving a file %s, %s,file will be resend", file.FileName(), err.Error())
				c.waiting.Store(file.FullPath(), file)
			} else {
				c.completed.Store(file.FullPath(), file)
			}
			c.inProgress.Delete(file.FullPath())

			c.log.Infof("sending %s completed", file.Metadata())
		case <-ctx.Done():
			return
		}
	}
}
func (c *Client) scan(ctx context.Context) {
	for {
		select {
		case <-time.After(time.Duration(c.opts.scanInterval) * time.Second):
			for _, folder := range c.absRootFolders {
				err := c.walk(folder)
				if err != nil {
					c.log.Errorf("error during scan files in folder %s, %s", folder, err.Error())
				}
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
