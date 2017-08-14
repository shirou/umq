package umq

// Message is a message
type Message struct {
	Body      []byte
	MessageID string
}

func (msg Message) GetMessageID() string {
	return msg.MessageID
}

func (msg Message) GetBody() []byte {
	return msg.Body
}
