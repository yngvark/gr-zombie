// Package integration implements functionality for more easily running integration tests
package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/rancher/k3d/v3/cmd/util"
)

// Localstack contains all state for managing the lifecycle
// of the localstack instance
type Localstack struct {
	pool     *dockertest.Pool
	edgePort int
	resource *dockertest.Resource
}

// NewLocalstack returns initialised localstack
// state
func NewLocalstack() *Localstack {
	return &Localstack{}
}

// Create a localstack instance
func (l *Localstack) Create(timeout time.Duration) error {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return fmt.Errorf("couldn't connect to docker: %w", err)
	}

	l.pool = pool

	port, err := util.GetFreePort()
	if err != nil {
		return fmt.Errorf("failed to find available port for edge: %w", err)
	}

	l.edgePort = port

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "localstack/localstack",
		Tag:        "0.11.4",
		Env: []string{
			"EDGE_PORT=4566",
			"DEFAULT_REGION=eu-west-1",
			"LAMBDA_REMOTE_DOCKER=0",
			"START_WEB=0",
			"DEBUG=1",
		},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"4566/tcp": {
				{
					HostIP:   "0.0.0.0",
					HostPort: fmt.Sprintf("%d", l.edgePort),
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to start localstack container: %w", err)
	}

	err = pool.Retry(l.Health)
	if err != nil {
		return fmt.Errorf("failed to wait for localstack: %w", err)
	}

	l.resource = resource

	return nil
}

// AWSSession returns an AWS session that works with the
// localstack instance
func (l *Localstack) AWSSession() *session.Session {
	return session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("eu-west-1"),
		Endpoint: aws.String(fmt.Sprintf("http://localhost:%d", l.edgePort)),
		// DisableSSL:  aws.Bool(true),
		Credentials: credentials.NewStaticCredentials("fake", "fake", "fake"),
	}))
}

// Logs retrieves the logs from the localstack instance
func (l *Localstack) Logs() (string, error) {
	var b bytes.Buffer

	err := l.pool.Client.Logs(docker.LogsOptions{
		Container:         l.resource.Container.ID,
		OutputStream:      &b,
		ErrorStream:       &b,
		InactivityTimeout: 0,
		Stdout:            true,
		Stderr:            true,
	})

	return b.String(), err
}

// Cleanup removes all resources created by localstack
func (l *Localstack) Cleanup() error {
	if l.pool != nil {
		err := l.pool.Purge(l.resource)
		if err != nil {
			return fmt.Errorf("failed to cleanup resources: %w", err)
		}
	}

	return nil
}

// Health returns nil if the localstack instance is up and running
func (l *Localstack) Health() error {
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/health?reload", l.edgePort))
	if err != nil {
		return err
	}

	var services map[string]map[string]string

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}

	err = json.Unmarshal(body, &services)
	if err != nil {
		return fmt.Errorf("failed to unmarshal json: %w", err)
	}

	for service, status := range services["services"] {
		if status != "running" {
			return fmt.Errorf("waiting for: %s, to get to running state, currently: %s", service, status)
		}
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got response code from localstack: %d, not 200 OK", resp.StatusCode)
	}

	return nil
}
