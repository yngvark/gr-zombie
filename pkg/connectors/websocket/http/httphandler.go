// Package http knows how to handle HTTP websocket connections
package http

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/gorilla/websocket"
	"github.com/yngvark/gridwalls3/source/zombie-go/pkg/pubsub"
)

// Handler knows how to handle HTTP websocket connections
type Handler struct {
	upgrader             *websocket.Upgrader
	connection           *websocket.Conn
	publisher            pubsub.Publisher
	stopGamelogicChannel chan bool
	log                  *zap.SugaredLogger
}

// NewHTTPHandler returns a new Handler
func NewHTTPHandler(logger *zap.SugaredLogger, allowedOrigins map[string]bool, publisher pubsub.Publisher, stopGamelogicChannel chan bool) *Handler {
	h := &Handler{
		log:                  logger,
		publisher:            publisher,
		stopGamelogicChannel: stopGamelogicChannel,
	}

	h.upgrader = &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			origin, ok := r.Header["Origin"]
			if !ok {
				return false
			}

			if len(origin) > 0 {
				_, ok := allowedOrigins[origin[0]]
				h.log.Infow("Checking origin %s. Result: %t\n", origin[0], ok)

				return ok
			}

			return true
		},
		EnableCompression: true,
	}

	return h
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	connection, err := h.upgrader.Upgrade(writer, request, nil)
	if err != nil {
		h.log.Error("could not upgrade:", err)
		return
	}

	// This only support one client.
	h.connection = connection

	h.log.Info("Client connected!")

	// Handle disconnection
	activelyCloseConnectionChannel := make(chan bool)

	defer h.closeConnectionWhenDone(activelyCloseConnectionChannel)

	go h.readIncomingMessages(activelyCloseConnectionChannel)
}

func (h *Handler) closeConnectionWhenDone(closeConnectionChannel chan bool) {
	select {
	case <-h.stopGamelogicChannel:
	case <-closeConnectionChannel:
	}

	h.log.Info("Closing connection from server")

	err := h.connection.Close()

	if err != nil {
		h.log.Info("error when closing connection: %w", err)
	} else {
		h.log.Info("Connection closed successfully.")
	}
}

func (h *Handler) readIncomingMessages(closeConnectionChannel chan bool) {
	for {
		h.log.Info("Reading next message...")

		_, message, err := h.connection.ReadMessage()
		if err != nil {
			// Client disconnected
			h.log.Info("Client disconnected")

			// We need to stop both game logic and disconnect
			h.stopGamelogicChannel <- true
			closeConnectionChannel <- true

			h.log.Errorf("Read error: %s", err.Error())

			return
		}

		err = h.handleIncomingMsg(message)
		if err != nil {
			h.log.Errorf("Error handling incoming message. Aborting. Error: %s", err.Error())
			return
		}
	}
}

func (h *Handler) handleIncomingMsg(message []byte) error {
	h.log.Infof("Received: %s", message)
	msgString := string(message)

	err := h.publisher.SendMsg(msgString)
	if err != nil {
		return fmt.Errorf("sending message with publisher: %w", err)
	}

	return nil
}

/*
func (h *Handler) sendMsg(msg string) error {
	if h.connection == nil {
		return errors.New("could not send message, not connected")
	}

	h.log.Infof("Sending msg: %s", msg)

	err := h.connection.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		return fmt.Errorf("could not write message: %w", err)
	}

	return nil
}
*/
