package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	awss3 "github.com/aws/aws-sdk-go/service/s3"
	"testing-stuff/pkg/env"
	"testing-stuff/pkg/img"
	"testing-stuff/platform/s3"
)

func main() {
	sess := session.Must(session.NewSession())
	awss3Client := awss3.New(sess)

	uploaderDownloader := s3.New(awss3Client)
	resizer := img.NewResizer()
	encoder := img.NewEncoder()

	handler := NewHandler(uploaderDownloader, resizer, encoder, env.Require("IMAGES_BUCKET_NAME"))

	lambda.Start(handler)
}
