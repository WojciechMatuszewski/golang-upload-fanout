// +build !test

package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	awss3 "github.com/aws/aws-sdk-go/service/s3"
	awssns "github.com/aws/aws-sdk-go/service/sns"
	"testing-stuff/pkg/env"
	"testing-stuff/pkg/formdata"
	"testing-stuff/platform/s3"
	"testing-stuff/platform/sns"
)

func main() {
	sess := session.Must(session.NewSession())

	reader := formdata.NewReader()
	uploader := s3.New(awss3.New(sess))
	publisher := sns.NewPublisher(awssns.New(sess), env.Require("SNS_TOPIC_ARN"))

	lambda.Start(NewHandler(reader, uploader, publisher, env.Require("IMAGES_BUCKET_NAME")))
}
