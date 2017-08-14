package umq

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type SQSTransport struct {
	svc *sqs.SQS
}

func NewSQSTransport(config *aws.Config) *SQSTransport {
	sess := session.Must(session.NewSession(config))
	svc := sqs.New(sess)

	return &SQSTransport{
		svc: svc,
	}
}

func (q *SQSTransport) Connect(url string) error {
	return nil
}

func (tr *SQSTransport) GetQueue(key string, opts ...Option) (Queue, error) {
	resultURL, err := tr.svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(key),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == sqs.ErrCodeQueueDoesNotExist {
			return nil, fmt.Errorf("queue not found: %v", key)
		}
		return nil, fmt.Errorf("unable to queue: %v, %v", key, err)
	}

	return &SQSQueue{
		svc:       tr.svc,
		QueueName: key,
		QueueURL:  resultURL,
	}, nil
}

func (tr *SQSTransport) CreateQueue(queue string) error {
	_, err := tr.svc.CreateQueue(&sqs.CreateQueueInput{
		QueueName: aws.String(queue),
	})
	return err
}
func (tr *SQSTransport) DeleteQueue(queue string) error {
	resultURL, err := tr.svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(queue),
	})
	if err != nil {
		return err
	}

	_, err = tr.svc.DeleteQueue(&sqs.DeleteQueueInput{
		QueueUrl: resultURL.QueueUrl,
	})
	return err
}

type SQSQueue struct {
	svc       *sqs.SQS
	QueueName string
	QueueURL  *sqs.GetQueueUrlOutput
}

func (q *SQSQueue) Close() error {
	return nil
}

func (q *SQSQueue) Receive(opts ...Option) (Message, error) {
	return q.ReceiveWithContext(context.Background(), opts...)
}

func (q *SQSQueue) ReceiveWithContext(ctx context.Context, opts ...Option) (Message, error) {
	timeout := int64(10)

	result, err := q.svc.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
		MaxNumberOfMessages: aws.Int64(1), // only One message is fetched
		QueueUrl:            q.QueueURL.QueueUrl,
		AttributeNames: aws.StringSlice([]string{
			"SentTimestamp",
		}),
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		WaitTimeSeconds: aws.Int64(timeout),
	})
	if err != nil {
		return Message{}, fmt.Errorf("unable to recieve message from queue: %s %v", q.QueueName, err)
	}

	if len(result.Messages) == 0 {
		return Message{}, fmt.Errorf("empty message recieved")
	}

	return q.sqsMsgToMessage(result.Messages[0]) // only need first message
}

func (q *SQSQueue) Delete(msg Message, opts ...Option) error {
	return q.DeleteWithContext(context.Background(), msg, opts...)
}

func (q *SQSQueue) DeleteWithContext(ctx context.Context, msg Message, opts ...Option) error {
	input := &sqs.DeleteMessageInput{
		QueueUrl:      q.QueueURL.QueueUrl,
		ReceiptHandle: aws.String(msg.GetMessageID()),
	}

	_, err := q.svc.DeleteMessageWithContext(ctx, input)
	return err
}

func (q *SQSQueue) Send(msg Message, opts ...Option) error {
	return q.SendWithContext(context.Background(), msg, opts...)
}

func (q *SQSQueue) SendWithContext(ctx context.Context, msg Message, opts ...Option) error {
	input := &sqs.SendMessageInput{
		//		DelaySeconds: aws.Int64(10),
		MessageBody: aws.String(string(msg.GetBody())),
		QueueUrl:    q.QueueURL.QueueUrl,
	}

	_, err := q.svc.SendMessageWithContext(ctx, input)
	return err
}

func (q *SQSQueue) sqsMsgToMessage(sm *sqs.Message) (Message, error) {
	msg := Message{
		Body:      []byte(*sm.Body),
		MessageID: *sm.MessageId,
	}

	return msg, nil
}
