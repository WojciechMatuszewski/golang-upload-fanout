package main_test

import (
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing-stuff/functions/upload"
	"testing-stuff/functions/upload/mock"
)

const mockBucketName string = "BUCKET_NAME"

func TestUpload(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("errors when reader returns err", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reader := mock.NewMockReader(ctrl)
		uploader := mock.NewMockUploader(ctrl)
		handler := main.NewHandler(reader, uploader, mockBucketName)

		reader.EXPECT().ReadBase64Encoded("FOO", "BAR").Return(nil, errors.New("an error"))

		res, err := handler(ctx, events.APIGatewayProxyRequest{Body: "FOO", Headers: map[string]string{"Content-Type": "BAR"}})

		assert.Equal(t, err.Error(), "an error")
		assert.Equal(t, res, events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		})
	})

	t.Run("fails when no Content-Type header is specified", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reader := mock.NewMockReader(ctrl)
		uploader := mock.NewMockUploader(ctrl)
		handler := main.NewHandler(reader, uploader, mockBucketName)

		res, err := handler(ctx, events.APIGatewayProxyRequest{})

		assert.NoError(t, err)
		assert.Equal(t, res, events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Content Type header not found",
		})
	})

	t.Run("reader succeeds but uploader fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reader := mock.NewMockReader(ctrl)
		uploader := mock.NewMockUploader(ctrl)
		handler := main.NewHandler(reader, uploader, mockBucketName)

		filePart := &multipart.FileHeader{
			Filename: "",
			Header:   nil,
			Size:     0,
		}

		reader.EXPECT().ReadBase64Encoded("FOO", "BAR").Return(&multipart.Form{
			File: map[string][]*multipart.FileHeader{
				"file": {
					filePart,
				},
			},
		}, nil)

		uploader.EXPECT().Upload(ctx, mockBucketName, filePart).Return("", errors.New("boom"))

		res, err := handler(ctx, events.APIGatewayProxyRequest{Body: "FOO", Headers: map[string]string{"Content-Type": "BAR"}})

		assert.Equal(t, err.Error(), "boom")
		assert.Equal(t, res, events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		})
	})

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reader := mock.NewMockReader(ctrl)
		uploader := mock.NewMockUploader(ctrl)
		handler := main.NewHandler(reader, uploader, mockBucketName)

		filePart := &multipart.FileHeader{
			Filename: "",
			Header:   nil,
			Size:     0,
		}
		secondFilePart := &multipart.FileHeader{
			Filename: "",
			Header:   nil,
			Size:     0,
		}

		reader.EXPECT().ReadBase64Encoded("FOO", "BAR").Return(&multipart.Form{
			File: map[string][]*multipart.FileHeader{
				"file": {
					filePart,
				},
				"secondFile": {
					secondFilePart,
				},
			},
		}, nil)

		firstFileLocation := "first_file"
		secondFileLocation := "second_file"
		uploader.EXPECT().Upload(ctx, mockBucketName, filePart).Return(firstFileLocation, nil)
		uploader.EXPECT().Upload(ctx, mockBucketName, filePart).Return(secondFileLocation, nil)

		res, err := handler(ctx, events.APIGatewayProxyRequest{Body: "FOO", Headers: map[string]string{"Content-Type": "BAR"}})

		assert.NoError(t, err)

		resB, err := json.Marshal([]string{firstFileLocation, secondFileLocation})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, res, events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       string(resB),
		})
	})

}
