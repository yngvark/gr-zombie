package kafka_connector

import (
	"github.com/oslokommune/okctl/pkg/cloud"
	"github.com/stretchr/testify/assert"
	"github.com/yngvark/gridwalls3/source/zombie-go/pkg/integration"
	"testing"
	"time"
)

func TestOnLocalstack(t *testing.T) {
	s := integration.NewLocalstack()
	err := s.Create(3 * time.Minute)
	assert.NoError(t, err)

	provider, err := cloud.NewFromSession("eu-west-1", "", s.AWSSession())
	assert.NoError(t, err)

	testCases := []struct {
		name string
	}{
		{
			name: "Should work",
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// create consumer, publisher, send message, verify message sent with content

		})
	}

}
