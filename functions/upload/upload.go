package main

import (
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"testing-stuff/platform/s3"
)

// Handler describes shape of a upload lambda function
type Handler func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

// Reader reads base64Encoded (from binary) form-data passed to lambda as request.Body
type Reader interface {
	ReadBase64Encoded(d, ct string) (*multipart.Form, error)
}

// Uploader uploads File to s3
type Uploader interface {
	UploadFromMultipart(ctx context.Context, bucket string, fo s3.FileOpener, fi s3.FileInfo) (string, error)
}

// Publisher publishes to SNS topic
type Publisher interface {
	Publish(ctx context.Context, payload string, attributes map[string]string) error
}

// NewHandler creates a handler for lambda which is responsible for uploading image to s3 and sending message to SNS topic
func NewHandler(reader Reader, uploader Uploader, publisher Publisher, bucket string) Handler {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

		ct, found := request.Headers["Content-Type"]
		if !found {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       "Content Type header not found",
			}, nil
		}

		f, err := reader.ReadBase64Encoded(request.Body, ct)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       err.Error(),
			}, err
		}

		var fLocations []string
		for _, file := range f.File {
			for _, filePart := range file {

				fi := s3.FileInfo{
					FileName:           filePart.Filename,
					ContentType:        filePart.Header.Get("Content-Type"),
					ContentDisposition: filePart.Header.Get("Content-Disposition"),
				}
				fLocation, err := uploader.UploadFromMultipart(ctx, bucket, filePart, fi)

				if err != nil {
					return events.APIGatewayProxyResponse{
						StatusCode: http.StatusBadRequest,
						Body:       err.Error(),
					}, err
				}

				mfi, err := json.Marshal(fi)
				if err != nil {
					return events.APIGatewayProxyResponse{
						StatusCode: http.StatusBadRequest,
						Body:       err.Error(),
					}, err
				}
				err = publisher.Publish(ctx, string(mfi), map[string]string{})
				if err != nil {
					return events.APIGatewayProxyResponse{
						StatusCode: http.StatusBadRequest,
						Body:       err.Error(),
					}, err
				}

				fLocations = append(fLocations, fLocation)
			}
		}

		err = f.RemoveAll()
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       err.Error(),
			}, err
		}

		jb, err := json.Marshal(fLocations)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       err.Error(),
			}, err
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string(jb),
		}, nil
	}
}
