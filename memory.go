package umq

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type MemoryTransport struct {
	queueMap map[string]*MemoryQueue
}

func NewMemoryTransport() *MemoryTransport {
	return &MemoryTransport{
		queueMap: make(map[string]*MemoryQueue),
	}
}

func (q *MemoryTransport) Connect(url string) error {
	return nil
}

func (tr *MemoryTransport) GetQueue(key string, opts ...Option) (Queue, error) {
	queue := &MemoryQueue{
		lock:      new(sync.RWMutex),
		IsConsume: true,
	}
	for _, opt := range opts {
		if !optionAvailable("memory", opt.Target()) {
			return nil, fmt.Errorf("option does not match for Memory, %v", opt)
		}
		if err := opt.Apply(queue); err != nil {
			return nil, err
		}
	}
	queue.QueueName = key
	tr.queueMap[key] = queue

	return queue, nil
}

type MemoryQueue struct {
	QueueName string
	queue     []Message
	IsConsume bool
	lock      *sync.RWMutex
}

func (q *MemoryQueue) Close() error {
	q.queue = make([]Message, 0)

	return nil
}

func (q *MemoryQueue) Receive(opts ...Option) (Message, error) {
	return q.ReceiveWithContext(context.Background(), opts...)
}

func (q *MemoryQueue) ReceiveWithContext(ctx context.Context, opts ...Option) (Message, error) {
	ch := make(chan Message)
	go func(ch chan Message) {
		c := time.Tick(1 * time.Millisecond)
		for range c {
			if len(q.queue) > 0 {
				q.lock.Lock()
				defer q.lock.Unlock()
				msg := q.queue[0]
				q.queue = q.queue[1:] // pop first message
				ch <- msg
				break
			}
		}
	}(ch)

	select {
	case <-ctx.Done():
		return Message{}, ctx.Err()
	case msg := <-ch:
		if q.IsConsume {
			if err := q.DeleteWithContext(ctx, msg); err != nil {
				return msg, err
			}
		}
		return msg, nil
	}

	return Message{}, nil
}

func (q *MemoryQueue) Delete(msg Message, opts ...Option) error {
	return q.DeleteWithContext(context.Background(), msg, opts...)
}

func (q *MemoryQueue) DeleteWithContext(ctx context.Context, msg Message, opts ...Option) error {
	q.lock.Lock()
	defer q.lock.Unlock()

	var index int
	for i, qm := range q.queue {
		if msg.GetMessageID() == qm.GetMessageID() {
			index = i
			break
		}
	}

	if index > 0 {
		q.queue = append(q.queue[:index-1], q.queue[index+1:]...)
	}

	return nil
}

func (q *MemoryQueue) Send(msg Message, opts ...Option) error {
	return q.SendWithContext(context.Background(), msg, opts...)
}

func (q *MemoryQueue) SendWithContext(ctx context.Context, msg Message, opts ...Option) error {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.queue = append(q.queue, msg)
	return nil
}
