package main

import (
	"altor/internal/settings"
	"altor/internal/websocket"
	log "github.com/sirupsen/logrus"
)

var cfg = settings.NewConfig()

func main() {
	settings.ConfigureLogger(cfg.LogLevel)
	log.Info("Starting Altor!")

	anilibriaWS, err := websocket.NewAnilibriaWS(cfg.AnilibriaWSURL.GetURL())
	if err != nil {
		log.Panicln(err)
	}

	err = anilibriaWS.ListenSubscribes()
	if err != nil {
		log.Panicln(err)
	}
}
