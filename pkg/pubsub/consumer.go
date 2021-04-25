package pubsub

type Subscriber interface {
	MsgReceived(msg string)
}

type Consumer interface {
	// ListenForMessages receives messages. It blocks until the Consumer's context is canceled.
	ListenForMessages()

	// Close closes the Consumer.
	Close() error
}
