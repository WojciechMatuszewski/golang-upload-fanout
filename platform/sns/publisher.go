package sns

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
)

// Publisher is a simple SNS client that enables publishing to a SNS topic
type Publisher struct {
	topicArn string
	sns      snsiface.SNSAPI
}

// NewPublisher is a function that creates Publisher
func NewPublisher(sns snsiface.SNSAPI, topicArn string) *Publisher {
	return &Publisher{
		topicArn: topicArn,
		sns:      sns,
	}
}

// Publish sends a message to given SNS topic
func (p *Publisher) Publish(ctx context.Context, payload string, attributes map[string]string) error {
	input := &sns.PublishInput{
		Message:  aws.String(payload),
		TopicArn: aws.String(p.topicArn),
	}

	if len(attributes) == 0 {
		_, err := p.sns.PublishWithContext(ctx, input)
		return err
	}

	msgAttr := make(map[string]*sns.MessageAttributeValue, len(attributes))
	for k, v := range attributes {
		msgAttr[k] = &sns.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(v),
		}

	}
	input.MessageAttributes = msgAttr

	_, err := p.sns.PublishWithContext(ctx, input)
	return err
}
