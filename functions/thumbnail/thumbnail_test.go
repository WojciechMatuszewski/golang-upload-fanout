package main_test

import (
	"bytes"
	"context"
	"errors"
	"image"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing-stuff/functions/thumbnail"
	"testing-stuff/functions/thumbnail/mock"
	"testing-stuff/platform/s3"
)

func TestThumbnail(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	const bucketName = "BUCKET_NAME"

	// why did I even mocked this ? wtf
	t.Run("encode failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		resizer := mock.NewMockThumbnailResizer(ctrl)
		uploaderDownloader := mock.NewMockUploaderDownloader(ctrl)
		encoder := mock.NewMockEncoder(ctrl)
		event := events.SQSEvent{Records: []events.SQSMessage{events.SQSMessage{
			Body: `{"Message": "{\"fileName\":\"filename.jpg\",\"contentType\":\"image/BAR\"}"}`,
		}}}
		img := prepareImage(t)
		handler := main.NewHandler(uploaderDownloader, resizer, encoder, bucketName)

		uploaderDownloader.EXPECT().Download(ctx, bucketName, "filename.jpg").Return(&img, nil)
		resizer.EXPECT().ResizeToThumbnail(img).Return(img)
		encoder.EXPECT().EncodeToBuffer(img, "image/BAR").Return(nil, errors.New("boom"))

		err := handler(ctx, event)

		assert.Error(t, err)
		assert.Equal(t, "boom", err.Error())
	})

	t.Run("Downloader failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		resizer := mock.NewMockThumbnailResizer(ctrl)
		uploaderDownloader := mock.NewMockUploaderDownloader(ctrl)
		encoder := mock.NewMockEncoder(ctrl)
		event := events.SQSEvent{Records: []events.SQSMessage{events.SQSMessage{
			Body: `{"Message": "{\"fileName\":\"filename.jpg\",\"contentType\":\"image/BAR\"}"}`,
		}}}
		handler := main.NewHandler(uploaderDownloader, resizer, encoder, bucketName)

		uploaderDownloader.EXPECT().Download(ctx, bucketName, "filename.jpg").Return(nil, errors.New("boom"))
		err := handler(ctx, event)

		assert.Error(t, err)
		assert.Equal(t, "boom", err.Error())
	})

	t.Run("upload resized failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		resizer := mock.NewMockThumbnailResizer(ctrl)
		uploaderDownloader := mock.NewMockUploaderDownloader(ctrl)
		encoder := mock.NewMockEncoder(ctrl)
		event := events.SQSEvent{Records: []events.SQSMessage{events.SQSMessage{
			Body: `{"Message": "{\"fileName\":\"filename.jpg\",\"contentType\":\"image/jpeg\"}"}`,
		}}}
		img := prepareImage(t)
		handler := main.NewHandler(uploaderDownloader, resizer, encoder, bucketName)
		buf := bytes.NewBuffer([]byte{})

		uploaderDownloader.EXPECT().Download(ctx, bucketName, "filename.jpg").Return(&img, nil)
		resizer.EXPECT().ResizeToThumbnail(img).Return(img)
		encoder.EXPECT().EncodeToBuffer(img, "image/jpeg").Return(buf, nil)
		uploaderDownloader.EXPECT().UploadFromReader(ctx, bucketName, buf, s3.FileInfo{
			FileName:    "thumb_filename.jpg",
			ContentType: "image/jpeg",
		}).Return("", errors.New("boom"))

		err := handler(ctx, event)

		assert.Error(t, err)
		assert.Equal(t, "boom", err.Error())
	})

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		resizer := mock.NewMockThumbnailResizer(ctrl)
		uploaderDownloader := mock.NewMockUploaderDownloader(ctrl)
		encoder := mock.NewMockEncoder(ctrl)
		event := events.SQSEvent{Records: []events.SQSMessage{events.SQSMessage{
			Body: `{"Message": "{\"fileName\":\"filename.jpg\",\"contentType\":\"image/jpeg\"}"}`,
		}}}
		img := prepareImage(t)
		handler := main.NewHandler(uploaderDownloader, resizer, encoder, bucketName)
		buf := bytes.NewBuffer([]byte{})

		uploaderDownloader.EXPECT().Download(ctx, bucketName, "filename.jpg").Return(&img, nil)
		resizer.EXPECT().ResizeToThumbnail(img).Return(img)
		encoder.EXPECT().EncodeToBuffer(img, "image/jpeg").Return(buf, nil)
		uploaderDownloader.EXPECT().UploadFromReader(ctx, bucketName, buf, s3.FileInfo{
			FileName:    "thumb_filename.jpg",
			ContentType: "image/jpeg",
		}).Return("file_location", nil)

		err := handler(ctx, event)

		assert.NoError(t, err)

	})
}

func prepareImage(t *testing.T) image.Image {
	t.Helper()

	rgba := image.NewNRGBA(image.Rect(1, 1, 2, 2))
	return rgba.SubImage(image.Rect(1, 1, 2, 2))
}
