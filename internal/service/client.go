package service

import (
	"context"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// DefaultDownloadPartSize is the default range of bytes to get at a time when
// using Download().
const DefaultDownloadPartSize = 1024 * 1024 * 5

// DefaultDownloadConcurrency is the default number of goroutines to spin up
// when using Download().

var DefaultDownloadConcurrency = 5

// AWS is a service to interact with original AWS services
type AWS interface {
	S3get(bucket, key string, rangeHeader *string) (*s3.GetObjectOutput, error)
	S3Download(f http.ResponseWriter, bucket, key string, rangeHeader *string) error
	S3listObjects(bucket, prefix string) (*s3.ListObjectsOutput, error)
	S3upload(bucket, key string, reader io.Reader) (output *s3manager.UploadOutput, err error)
}

type client struct {
	context.Context
	// *session.Session
	S3          s3iface.S3API
	session     *session.Session
	PartSize    int64
	Concurrency int
}

// NewClient returns new AWS client
func NewClient(ctx context.Context, region *string, partSize int64, concurrency int) AWS {
	sess := awsSession(region)
	return &client{Context: ctx,
		session:     sess,
		S3:          s3.New(sess),
		PartSize:    partSize,
		Concurrency: concurrency,
	}
}
