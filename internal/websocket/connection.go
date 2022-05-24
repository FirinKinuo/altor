package websocket

import (
	"altor/internal/websocket/messages"
	"altor/internal/websocket/messages/connect"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type AnilibriaWS struct {
	conn *websocket.Conn
}

// NewAnilibriaWS - constructor of the AnilibriaWS struct, accepts the URL of the Anilibria websocket as input
func NewAnilibriaWS(wsURL string) (ws *AnilibriaWS, err error) {
	connection, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return nil, err
	}

	return &AnilibriaWS{
		conn: connection,
	}, nil
}

// CloseConnection with Anilibria websocket
// If an attempt to close the connection causes an error, log.Panic is called
func (w *AnilibriaWS) CloseConnection() {
	err := w.conn.Close()
	if err != nil {
		log.Panic(err)
	}
}

// CheckConnection with Anilibria WebSocket
func (w *AnilibriaWS) CheckConnection(response []byte) (err error) {
	connectMessage, err := connect.Parse(response)
	if err != nil {
		return err
	}

	err = connectMessage.CheckConnection()
	if err != nil {
		return err
	}

	log.Infof("Connection to Anilibria Websocket success")
	return nil
}

// ParseMessage from WebSocket response
// Depending on the response received, the corresponding methods will be executed.
func (w *AnilibriaWS) ParseMessage(response []byte) (err error) {
	parsedMessageType := messages.ParseMessageType(response)

	switch parsedMessageType.Type {
	case messages.Connection:
		err = w.CheckConnection(response)
		if err != nil {
			log.Panic(err)
		}
	}

	return nil
}

// ListenSubscribes - starts the main loop listening for a websocket connection
// When a message is received, calls the parser ParseMessage
func (w *AnilibriaWS) ListenSubscribes() (err error) {
	defer w.CloseConnection()

	for {
		_, message, err := w.conn.ReadMessage()
		if err != nil {
			return err
		}

		err = w.ParseMessage(message)
		if err != nil {
			log.Error(err)
		}
	}
}
