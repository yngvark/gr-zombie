package pubsub_test

import (
	"testing"

	net "github.com/yngvark/gridwalls3/source/zombie-go/pkg/pubsub"

	"github.com/stretchr/testify/assert"
)

func TestPubSub(t *testing.T) {
	t.Run("Should send message to listeners", func(t *testing.T) {
		// Given
		b := net.NewBroker()
		testReceiver := &testReceiver{}

		b.Subscribe(testReceiver)

		// When
		b.SendMsg("YO")

		// Then
		assert.Equal(t, "YO", testReceiver.lastMsgReceived)
	})
}

type testReceiver struct {
	lastMsgReceived string
}

func (t *testReceiver) MsgReceived(msg string) {
	t.lastMsgReceived = msg
}
