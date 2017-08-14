package umq

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
)

func TestSQSQueue(t *testing.T) {
	assert := assert.New(t)
	tr := NewSQSTransport(&aws.Config{
		Endpoint: aws.String("http://localhost:9324"),
		Region:   aws.String("us-east-1"),
	})
	assert.Nil(tr.Connect("0.0.0.0:9324"))
	queueName := "A"
	assert.Nil(tr.CreateQueue(queueName))

	t.Run("create queue", func(t *testing.T) {
		time.Sleep(1)
		q, err := tr.GetQueue(queueName)
		assert.Nil(err)
		assert.NotNil(q)
	})

	t.Run("read/write", func(t *testing.T) {
		q, err := tr.GetQueue(queueName)
		assert.Nil(err)
		body := []byte("A")
		q.Send(Message{Body: body})
		msg, err := q.Receive()
		assert.Nil(err)
		assert.Equal(body, msg.GetBody())
	})
	t.Run("blocked", func(t *testing.T) {
		q, err := tr.GetQueue(queueName)
		assert.Nil(err)
		body := []byte("A")
		var wg sync.WaitGroup
		go func() {
			wg.Add(1)
			defer wg.Done()
			msg, err := q.Receive()
			assert.Nil(err)
			assert.Equal(body, msg.GetBody())
		}()
		time.Sleep(1 * time.Second)
		q.Send(Message{Body: body})
		wg.Wait()
	})
	t.Run("consume_option", func(t *testing.T) {
		q, err := tr.GetQueue(queueName,
			NewConsumeOption(false),
		)
		assert.Nil(err)
		body := []byte("A")
		var wg sync.WaitGroup
		var mid string
		go func() {
			wg.Add(1)
			defer wg.Done()

			msg, err := q.Receive()
			assert.Nil(err)
			assert.Equal(body, msg.GetBody())
			mid = msg.GetMessageID()
		}()
		time.Sleep(1 * time.Second)
		q.Send(Message{Body: body})
		wg.Wait()
		assert.Nil(q.Delete(Message{MessageID: mid}))
	})
	t.Run("context_cancel", func(t *testing.T) {
		q, err := tr.GetQueue(queueName)
		assert.Nil(err)
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
		msg, err := q.ReceiveWithContext(ctx)
		assert.NotNil(err)
		assert.Empty(msg.GetBody())
	})

	assert.Nil(tr.DeleteQueue(queueName))
}
