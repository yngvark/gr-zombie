package pubsub

type Publisher interface {
	// SendMsg sends messages
	SendMsg(msg string) error

	// Close closes the publisher
	Close()
}
