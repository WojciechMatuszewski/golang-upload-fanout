package s3

import (
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3 struct {
	s3Uploader   *s3manager.Uploader
	s3Downloader *s3manager.Downloader
}

type FileInfo struct {
	FileName           string `json:"fileName"`
	ContentType        string `json:"contentType"`
	ContentDisposition string `json:"contentDisposition"`
}

func New(s3client s3iface.S3API) *S3 {
	return &S3{
		s3Uploader:   s3manager.NewUploaderWithClient(s3client),
		s3Downloader: s3manager.NewDownloaderWithClient(s3client),
	}
}
