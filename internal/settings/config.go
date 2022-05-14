package settings

import (
	"flag"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Debug          bool
	LogLevel       string
	AnilibriaWSURL string
}

func NewConfig() (cfg *Config) {
	debug := flag.Bool("debug", false, "enable debug mode")
	logLevel := flag.String("log-level", "info", "set logging level")
	anilibriaWSURL := flag.String(
		"anilibria-ws-host",
		"wss://api.anilibria.tv/v2/ws/",
		"URL Anilibria websocket")

	flag.Parse()

	// Switching the logger level to debug mode
	// It is logical that during debugging all logs are needed
	if *debug {
		*logLevel = "debug"
	}

	return &Config{
		Debug:          *debug,
		LogLevel:       *logLevel,
		AnilibriaWSURL: *anilibriaWSURL,
	}
}

func ConfigureLogger(logLevel string) {
	log.StandardLogger().SetFormatter(&log.JSONFormatter{TimestampFormat: "2006-01-02T15:04:05.999"})
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		level = log.InfoLevel
		log.WithError(err).
			WithField("event", "start.parse_level_fail").
			Info("set log level to \"info\" by default")
	}
	log.SetLevel(level)
}
