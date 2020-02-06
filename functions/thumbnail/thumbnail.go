package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"io"

	"github.com/aws/aws-lambda-go/events"
	"testing-stuff/platform/s3"
)

type SNSMessage struct {
	Message string `json:"Message"`
}

type Handler func(ctx context.Context, event events.SQSEvent) error

// ThumbnailResizer resizes image to a thumbnail
type ThumbnailResizer interface {
	ResizeToThumbnail(img image.Image) *image.NRGBA
}

// UploaderDownloader downloads files from s3 and uploads them to s3 in form of io.Readers
type UploaderDownloader interface {
	UploadFromReader(ctx context.Context, bucket string, reader io.Reader, fi s3.FileInfo) (string, error)
	Download(ctx context.Context, bucket, key string) (*image.Image, error)
}

// Encoder encodes resized image to buffer
type Encoder interface {
	EncodeToBuffer(img image.Image, ext string) (*bytes.Buffer, error)
}

func NewHandler(uploaderDownloader UploaderDownloader, resizer ThumbnailResizer, encoder Encoder, bucket string) Handler {
	return func(ctx context.Context, event events.SQSEvent) error {
		message := event.Records[0]

		var snsm SNSMessage
		err := json.Unmarshal([]byte(message.Body), &snsm)
		if err != nil {
			fmt.Println("error while unmarshalling", err)
			return err
		}

		var s3fi s3.FileInfo
		err = json.Unmarshal([]byte(snsm.Message), &s3fi)
		if err != nil {
			fmt.Println("error on second unmarshal pass")
			return err
		}

		img, err := uploaderDownloader.Download(ctx, bucket, s3fi.FileName)
		if err != nil {
			fmt.Println("error while downloading", err)
			return err
		}

		resized := resizer.ResizeToThumbnail(*img)

		b, err := encoder.EncodeToBuffer(resized, s3fi.ContentType)
		if err != nil {
			fmt.Println("error while encoding to buffer", err)
			return err
		}

		fl, err := uploaderDownloader.UploadFromReader(ctx, bucket, b, s3.FileInfo{
			FileName:    "thumb_" + s3fi.FileName,
			ContentType: s3fi.ContentType,
		})
		if err != nil {
			fmt.Println("error while uploading", err)
			return err
		}

		fmt.Println(fl)

		return nil
	}
}
