package umq

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRedisQueue(t *testing.T) {
	assert := assert.New(t)
	tr := NewRedisTransport()
	assert.Nil(tr.Connect("redis://localhost:6379/"))
	queueName := "A"

	t.Run("create queue", func(t *testing.T) {
		time.Sleep(1)
		q, err := tr.GetQueue(queueName)
		assert.Nil(err)
		assert.NotNil(q)
	})
	/* redis does not store message in a queue
	t.Run("read/write", func(t *testing.T) {
	*/
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
	/* redis can not delete messagequeue
	t.Run("consume_option", func(t *testing.T) {
	*/
	t.Run("context_cancel", func(t *testing.T) {
		q, err := tr.GetQueue(queueName)
		assert.Nil(err)
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
		msg, err := q.ReceiveWithContext(ctx)
		assert.NotNil(err)
		assert.Empty(msg.GetBody())
	})
}
