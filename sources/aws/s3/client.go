package s3

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemq-sources/config"
	"github.com/kubemq-io/kubemq-sources/middleware"
	"github.com/kubemq-io/kubemq-sources/pkg/logger"
	"sync"
	"time"
)

type Client struct {
	log         *logger.Logger
	opts        options
	client      *s3.S3
	downloader  *s3manager.Downloader
	waiting     sync.Map
	inProgress  sync.Map
	completed   sync.Map
	sendCh      chan *SourceFile
	ctx         context.Context
	cancelFunc  context.CancelFunc
	scanFolders map[string]string
}

func New() *Client {
	return &Client{
		waiting:     sync.Map{},
		inProgress:  sync.Map{},
		completed:   sync.Map{},
		scanFolders: map[string]string{},
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
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(c.opts.region),
		Credentials: credentials.NewStaticCredentials(c.opts.awsKey, c.opts.awsSecretKey, c.opts.token),
	})
	if err != nil {
		return err
	}

	svc := s3.New(sess)
	c.client = svc
	c.downloader = s3manager.NewDownloader(sess)
	for _, folder := range c.opts.folders {
		unixFolder := unixNormalize(folder)
		c.scanFolders[unixFolder] = unixFolder
	}
	buckets, err := c.client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return err
	}
	for _, bucket := range buckets.Buckets {
		if *bucket.Name == c.opts.bucketName {
			return nil
		}
	}

	return fmt.Errorf("scanned bucket name not found")
}

func (c *Client) Start(ctx context.Context, target middleware.Middleware) error {
	c.ctx, c.cancelFunc = context.WithCancel(ctx)
	c.sendCh = make(chan *SourceFile)
	go c.scan(c.ctx)
	for i := 0; i < c.opts.concurrency; i++ {
		go c.senderFunc(c.ctx, target)
	}

	go c.send(c.ctx)
	return nil
}

//
func (c *Client) inPipe(ctx context.Context, file *SourceFile) bool {
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
			if err := file.Delete(ctx); err != nil {
				c.log.Errorf("error during delete a file %s,%s, will try again", file.FullPath(), err.Error())
			}
			return true
		} else {
			c.log.Infof("file %s already sent but a new content has been detected, resending", file.FullPath())
		}
	}
	return false
}

//
func (c *Client) walk(ctx context.Context) error {
	var list []*SourceFile
	objList, err := c.client.ListObjectsV2WithContext(ctx, &s3.ListObjectsV2Input{Bucket: aws.String(c.opts.bucketName)})
	if err != nil {
		return err
	}
	var objects []*s3.Object
	objects = append(objects, objList.Contents...)

	for _, object := range objects {
		srcFile := NewSourceFile(c.client, c.downloader, c.opts.bucketName, object)
		_, ok := c.scanFolders["/"]
		if ok {
			list = append(list, srcFile)
			continue
		}
		_, ok = c.scanFolders[srcFile.RootDir()]
		if ok {
			list = append(list, srcFile)
		}
	}
	added := 0
	for _, file := range list {
		if !c.inPipe(ctx, file) {
			c.waiting.Store(file.FullPath(), file)
			added++
		}
	}
	if added > 0 {
		c.log.Debugf("%d new files added to sending waiting list", added)
	}
	return nil
}

//
func (c *Client) senderFunc(ctx context.Context, sender middleware.Middleware) {
	for {
		select {
		case file := <-c.sendCh:
			c.inProgress.Store(file.FullPath(), file)
			c.waiting.Delete(file.FullPath())
			req, err := file.Request(ctx, c.opts.targetType, c.opts.bucketName)
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
			if err := file.Do(ctx); err != nil {
				c.log.Errorf("error during delete/moving a file %s, %s,file will be resend", file.FileName(), err.Error())
				c.waiting.Store(file.FullPath(), file)
			} else {
				c.completed.Store(file.FullPath(), file)
			}
			c.inProgress.Delete(file.FullPath())

			c.log.Infof("sending %s completed: ", file.Metadata())
		case <-ctx.Done():
			return
		}
	}
}
func (c *Client) scan(ctx context.Context) {
	for {
		select {
		case <-time.After(time.Duration(c.opts.scanInterval) * time.Second):
			err := c.walk(ctx)
			if err != nil {
				c.log.Errorf("error during scan files in bucket %s, %s", c.opts.bucketName, err.Error())
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

func (c *Client) Stop() error {
	c.cancelFunc()
	c.waiting = sync.Map{}
	return nil
}
