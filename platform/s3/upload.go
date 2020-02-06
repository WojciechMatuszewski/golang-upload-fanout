package s3

import (
	"context"
	"errors"
	"io"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// FileOpener is an interface which is passed to Upload receiver
type FileOpener interface {
	Open() (multipart.File, error)
}

// UploadFromReader uploads images from io.Reader to s3
func (s *S3) UploadFromReader(ctx context.Context, bucket string, r io.Reader, fi FileInfo) (string, error) {

	result, err := s.s3Uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Body:               r,
		Bucket:             aws.String(bucket),
		Key:                aws.String(fi.FileName),
		ContentDisposition: aws.String(fi.ContentDisposition),
		ContentType:        aws.String(fi.ContentType),
	})
	if err != nil {
		return "", err
	}

	return result.Location, nil
}

// UploadFromMultipart uploads images to s3 from multipart requests
func (s *S3) UploadFromMultipart(ctx context.Context, bucket string, fo FileOpener, fi FileInfo) (string, error) {
	f, err := fo.Open()
	if err != nil {
		return "", errors.New("could not open the file")
	}

	result, err := s.s3Uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Body:               f,
		Bucket:             aws.String(bucket),
		Key:                aws.String(fi.FileName),
		ContentDisposition: aws.String(fi.ContentDisposition),
		ContentType:        aws.String(fi.ContentType),
	})
	if err != nil {
		return "", err
	}

	err = f.Close()
	if err != nil {
		return "", err
	}

	return result.Location, nil
}
