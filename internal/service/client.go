package service

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

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
	MinioUpload(bucketName, objectName, filePath string) (output minio.UploadInfo, err error)
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
	// 删除endpoint的协议
	endpoint := *sess.Config.Endpoint
	if strings.Contains(endpoint, "http://") {
		endpoint = strings.Trim(endpoint, "http://")
	} else if strings.Contains(endpoint, "https://") {
		endpoint = strings.Trim(endpoint, "https://")
	}
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
		Secure: true,
	})
	if err != nil {
		log.Fatalf("create minioClient fail, err:%s", err.Error())
	}
	return &client{Context: ctx,
		session:     sess,
		S3:          s3.New(sess),
		PartSize:    partSize,
		Concurrency: concurrency,
		MinioClient: minioClient,
	}
}
