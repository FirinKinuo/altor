package main

import (
	"altor/internal/settings"
	log "github.com/sirupsen/logrus"
)

var cfg = settings.NewConfig()

func main() {
	settings.ConfigureLogger(cfg.LogLevel)
	log.Info("Starting Altor!")
}
