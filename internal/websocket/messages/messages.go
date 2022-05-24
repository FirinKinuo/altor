package messages

import (
	"encoding/json"
)

type Message string

const (
	TorrentUpdate Message = "torrent_update"
	Subscribe             = "subscribe"
	Connection            = ""
)

type MessageType struct {
	Type Message `json:"type"`
}

func ParseMessageType(response []byte) (mt *MessageType) {
	var messageType_ MessageType

	_ = json.Unmarshal(response, &messageType_)

	return &messageType_
}
