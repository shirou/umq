package umq

import (
	"context"
	"time"

	redis "github.com/garyburd/redigo/redis"
)

type RedisTransport struct {
	pool *redis.Pool
}

func NewRedisTransport() *RedisTransport {
	return &RedisTransport{}
}

func (tr *RedisTransport) Connect(url string) error {
	tr.pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(url)
		},
	}

	return nil
}

func (q *RedisTransport) Close() error {
	return q.pool.Close()
}

func (tr *RedisTransport) GetQueue(key string, opts ...Option) (Queue, error) {
	return &RedisQueue{
		QueueName: key,
		IsConsume: true,
		transport: tr,
	}, nil
}

type RedisQueue struct {
	QueueName string
	IsConsume bool
	transport *RedisTransport
}

func (q *RedisQueue) Close() error {
	return nil
}

func (q *RedisQueue) Receive(opts ...Option) (Message, error) {
	return q.ReceiveWithContext(context.Background(), opts...)
}

func (q *RedisQueue) ReceiveWithContext(ctx context.Context, opts ...Option) (Message, error) {
	ch := make(chan Message)
	var retErr error
	go func(ch chan Message) {
		conn := q.transport.pool.Get()
		psc := redis.PubSubConn{Conn: conn}
		defer psc.Close()
		psc.Subscribe(q.QueueName)
		defer close(ch)
		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				ch <- Message{Body: v.Data}
				return
			case redis.Subscription:
				// do nothing
			case error:
				retErr = v
				ch <- Message{}
				return
			}
		}
	}(ch)

	select {
	case <-ctx.Done():
		return Message{}, ctx.Err()
	case msg := <-ch:
		return msg, retErr
	}

	return Message{}, retErr
}

func (q *RedisQueue) Delete(msg Message, opts ...Option) error {
	return q.DeleteWithContext(context.Background(), msg, opts...)
}

func (q *RedisQueue) DeleteWithContext(ctx context.Context, msg Message, opts ...Option) error {
	return nil
}

func (q *RedisQueue) Send(msg Message, opts ...Option) error {
	return q.SendWithContext(context.Background(), msg, opts...)
}

func (q *RedisQueue) SendWithContext(ctx context.Context, msg Message, opts ...Option) error {
	conn := q.transport.pool.Get()
	defer conn.Close()

	ch := make(chan bool)
	var retErr error
	go func(ch chan bool) {
		conn.Do("PUBLISH", q.QueueName, msg.Body)
		conn.Flush()
		ch <- true
	}(ch)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ch:
		return retErr
	}

	return retErr
}
