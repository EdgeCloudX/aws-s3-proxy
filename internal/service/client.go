package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/minio/minio-go/v7"
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
	S3Header(bucket, key string) (output *s3.HeadObjectOutput, err error)
}

type client struct {
	context.Context
	// *session.Session
	S3          s3iface.S3API
	session     *session.Session
	PartSize    int64
	Concurrency int
	MinioClient *minio.Client
}

// NewClient returns new AWS client
func NewClient(ctx context.Context, region *string, partSize int64, concurrency int) AWS {
	sess := awsSession(region)
	endPoint := "https://" + *sess.Config.Endpoint
	command := fmt.Sprintf("mc alias set s3 %s %s %s --api s3v4", endPoint, os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"))
	log.Printf("mc init Command: %s", command)
	// 需要执行的命令： free -mh
	cmd := exec.Command("/bin/sh", "-c", command)

	result, err := cmd.Output()
	if err != nil {
		log.Printf("mc init Command error: %s", err.Error())
	}
	log.Printf("shell result:%s", string(result))
	return &client{Context: ctx,
		session:     sess,
		S3:          s3.New(sess),
		PartSize:    partSize,
		Concurrency: concurrency,
	}
}
