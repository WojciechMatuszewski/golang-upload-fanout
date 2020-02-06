package sns_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	awssns "github.com/aws/aws-sdk-go/service/sns"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing-stuff/platform/sns"
	"testing-stuff/platform/sns/mock"
)

const topicArn = "TOPIC_ARN"

func TestPublisher_Publish(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("publishes with attributes", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		snsAPI := mock.NewMockSNSAPI(ctrl)
		publisher := sns.NewPublisher(snsAPI, topicArn)
		expectedInput := &awssns.PublishInput{
			Message:  aws.String("PAYLOAD"),
			TopicArn: aws.String(topicArn),
			MessageAttributes: map[string]*awssns.MessageAttributeValue{
				"key": {
					DataType:    aws.String("String"),
					StringValue: aws.String("attr1"),
				},
			},
		}

		snsAPI.EXPECT().PublishWithContext(ctx, expectedInput).Return(&awssns.PublishOutput{}, nil)
		err := publisher.Publish(ctx, "PAYLOAD", map[string]string{"key": "attr1"})

		assert.NoError(t, err)
	})

	t.Run("publishes with empty attributes", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		snsAPI := mock.NewMockSNSAPI(ctrl)
		publisher := sns.NewPublisher(snsAPI, topicArn)
		expectedInput := &awssns.PublishInput{
			Message:           aws.String("PAYLOAD"),
			TopicArn:          aws.String(topicArn),
			MessageAttributes: nil,
		}

		snsAPI.EXPECT().PublishWithContext(ctx, expectedInput).Return(&awssns.PublishOutput{}, nil)
		err := publisher.Publish(ctx, "PAYLOAD", map[string]string{})

		assert.NoError(t, err)
	})

	t.Run("publish failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		snsAPI := mock.NewMockSNSAPI(ctrl)
		publisher := sns.NewPublisher(snsAPI, topicArn)
		expectedInput := &awssns.PublishInput{
			Message:           aws.String("PAYLOAD"),
			TopicArn:          aws.String(topicArn),
			MessageAttributes: nil,
		}

		snsAPI.EXPECT().PublishWithContext(ctx, expectedInput).Return(nil, errors.New("boom"))
		err := publisher.Publish(ctx, "PAYLOAD", map[string]string{})

		assert.Error(t, err)
	})
}
