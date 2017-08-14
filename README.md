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
	assert.Nil(tr.Connect("0.0.0.0:9324"))

    // get queue by queue name
	q, err := tr.GetQueue(queueName)

    // send message.
    q.Send(Message{Body: body})

    // recieve. this function blocks.
    msg, err := q.Receive()
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
