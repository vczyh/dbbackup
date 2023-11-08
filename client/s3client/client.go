package s3client

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"sync"
)

type Client struct {
	accessKeyID     string
	secretAccessKey string
	forcePathStyle  bool
	region          string
	endpoint        string
	delimiter       string
	mu              sync.Mutex
	// use AwsClient() instead.
	_awsClient *s3.S3
}

func New(accessKeyID, secretAccessKey, endpoint string, opts ...Option) *Client {
	c := &Client{
		accessKeyID:     accessKeyID,
		secretAccessKey: secretAccessKey,
		endpoint:        endpoint,
	}
	for _, opt := range opts {
		opt.apply(c)
	}

	if c.delimiter == "" {
		c.delimiter = "/"
	}

	return c
}

func (c *Client) ListBuckets(ctx context.Context) ([]string, error) {
	awsClient, err := c.AwsClient()
	if err != nil {
		return nil, err
	}
	out, err := awsClient.ListBucketsWithContext(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}
	var buckets []string
	for _, bucket := range out.Buckets {
		buckets = append(buckets, *bucket.Name)
	}

	return buckets, nil
}

func (c *Client) AwsClient() (*s3.S3, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c._awsClient != nil {
		return c._awsClient, nil
	}

	// TODO tls self sign certification ignore?

	awsConfig := &aws.Config{
		// TODO  log level
		//LogLevel: aws.LogOff,
		Region:           &c.region,
		Endpoint:         &c.endpoint,
		S3ForcePathStyle: &c.forcePathStyle,
		Credentials:      credentials.NewStaticCredentials(c.accessKeyID, c.secretAccessKey, ""),
	}
	awsSession, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, err
	}
	c._awsClient = s3.New(awsSession)

	return c._awsClient, nil
}

type Option interface {
	apply(*Client)
}

type optionFunc func(*Client)

func (f optionFunc) apply(c *Client) {
	f(c)
}

func WithForcePathStyle(forcePathStyle bool) Option {
	return optionFunc(func(c *Client) {
		c.forcePathStyle = forcePathStyle
	})
}

func WithRegion(region string) Option {
	return optionFunc(func(c *Client) {
		c.region = region
	})
}
