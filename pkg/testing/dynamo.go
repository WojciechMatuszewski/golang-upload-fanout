package testing

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const dynamodbURL = "http://localhost:6000"

func localSession() *session.Session {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("local"),
		Credentials: credentials.NewStaticCredentials("local", "local", "local"),
	})

	if err != nil {
		panic(err)
	}
	return sess
}

func localDynamo() *dynamodb.DynamoDB {
	return dynamodb.New(localSession(), &aws.Config{Endpoint: aws.String(dynamodbURL)})
}

func SetupDynamoTest(t *testing.T) (context.Context, *dynamodb.DynamoDB) {
	t.Helper()

	return context.Background(), localDynamo()
}
