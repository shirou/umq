# umq (Universal Message Queue)

same as kombu


# usage


```
// create transport
tr := NewSQSTransport(&aws.Config{
Endpoint: aws.String("http://localhost:9324"),
Region:   aws.String("us-east-1"),
})

// connect transport
tr.Connect()

// get queue by queue name
q, err := tr.GetQueue(queueName)

// send message.
q.Send(Message{Body: body})

// recieve. this function blocks.
msg, err := q.Receive()

// If you want to timeout or cancel, use context
ctx := context.WithTimeout(context.Background())
q.SendWithContext(ctx, Message{Body: body})
msg, err := q.Receive(ctx)

```

Message is an abstract message which has body, `[]byte`.

```
type Message struct {
Body      []byte
MessageID string
}
```


# available transport

- memory
- Amazon SQS
- Redis Pub/Sub



# License

MIT
