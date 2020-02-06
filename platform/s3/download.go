package s3

import (
	"bytes"
	"context"
	"image"
	// These are needed to initialize different encoders, since downloaded image can be of different extension
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (s *S3) Download(ctx context.Context, bucket, key string) (*image.Image, error) {
	buf := aws.NewWriteAtBuffer([]byte{})

	_, err := s.s3Downloader.DownloadWithContext(ctx, buf, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	imgR := bytes.NewReader(buf.Bytes())

	img, _, err := image.Decode(imgR)
	if err != nil {
		return nil, err
	}

	return &img, nil

}
