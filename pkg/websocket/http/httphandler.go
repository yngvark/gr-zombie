package http

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/yngvark/gridwalls3/source/zombie-go/pkg/pubsub"
)

type HTTPHandler struct {
	upgrader             *websocket.Upgrader
	connection           *websocket.Conn
	publisher            pubsub.Publisher
	stopGamelogicChannel chan bool
	log                  *zap.SugaredLogger
}

func NewHTTPHandler(logger *zap.SugaredLogger, allowedOrigins map[string]bool, publisher pubsub.Publisher, stopGamelogicChannel chan bool) *HTTPHandler {
	h := &HTTPHandler{
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

func (h *HTTPHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
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

	defer h.CloseConnectionWhenDone(activelyCloseConnectionChannel)

	go h.ReadIncomingMessages(activelyCloseConnectionChannel)
}

func (h *HTTPHandler) CloseConnectionWhenDone(closeConnectionChannel chan bool) {
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

func (h *HTTPHandler) ReadIncomingMessages(closeConnectionChannel chan bool) {
	for {
		h.log.Info("Reading next message...")

		_, message, err := h.connection.ReadMessage()
		if err != nil {
			// Client disconnected
			h.log.Info("Client disconnected")

			// We need to stop both game logic and disconnect
			h.stopGamelogicChannel <- true
			closeConnectionChannel <- true

			h.log.Errorf("Read error: %w")

			break
		}

		h.HandleIncomingMsg(message)
	}
}

func (h *HTTPHandler) HandleIncomingMsg(message []byte) {
	h.log.Infof("Received: %s", message)
	msgString := string(message)
	h.publisher.SendMsg(msgString)
}

func (h *HTTPHandler) SendMsg(msg string) error {
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
