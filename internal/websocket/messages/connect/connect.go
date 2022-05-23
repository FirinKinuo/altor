package connect

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type (
	Response struct {
		Connection string `json:"connection"`
		ApiVersion string `json:"api_version"`
	}

	Message struct {
		Connected  bool
		ApiVersion string
	}
)

// Parse connection status message
func Parse(response []byte) (r *Message, err error) {
	var jsonResponse Response

	err = json.Unmarshal(response, &jsonResponse)
	if err != nil {
		return nil, fmt.Errorf("unable to read json Connection, reason: %w", err)
	}

	connected := false

	if strings.ToLower(jsonResponse.Connection) == "success" {
		connected = true
	}

	return &Message{
		Connected:  connected,
		ApiVersion: jsonResponse.ApiVersion,
	}, nil
}

// CheckConnection to the Anilibria WS server
func (m *Message) CheckConnection() (err error) {
	if !m.Connected {
		return errors.New("connection to Anilibria Websocket was not established")
	}
	return err
}
