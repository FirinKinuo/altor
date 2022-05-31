package settings

import (
	"flag"
	log "github.com/sirupsen/logrus"
)

// Flags structure with all arguments from the command line
type Flags struct {
	Debug                      bool
	LogLevel                   string
	AnilibriaWSScheme          string
	AnilibriaWSHost            string
	AnilibriaWSPath            string
	QBittorrentScheme          string
	QBittorrentHost            string
	QBittorrentIgnoreTSLVerify bool
	QBittorrentUser            string
	QBittorrentPassword        string
	QBittorrentSaveFolder      string
	QBittorrentCategory        string
}

// NewFlags - constructor of the Flags structure
func NewFlags() (flags *Flags) {
	debug := flag.Bool("debug", false, "enable debug mode")
	logLevel := flag.String("log-level", "info", "set logging level")
	anilibriaWSScheme := flag.String("anilibria-ws-scheme", "wss", "Scheme Anilibria websocket")
	anilibriaWSHost := flag.String("anilibria-ws-host", "api.anilibria.tv", "Host Anilibria websocket")
	anilibriaWSPath := flag.String("anilibria-ws-path", "v2/ws/", "Path Anilibria websocket")
	qBittorrentScheme := flag.String("qbt-scheme", "http", "Scheme of qBittorrent Web Client")
	qBittorrentHost := flag.String("qbt-host", "localhost:8080", "Host of qBittorrent Web Client")
	qBittorrentUser := flag.String("qbt-user", "admin", "User for qBittorrent Web Client")
	qBittorrentIgnoreTSLVerify := flag.Bool(
		"qbt-ignore-tls-verify",
		false,
		"Ignore TLS verify for self-signed certificate for qBittorrent Web Client")
	qBittorrentPassword := flag.String("qbt-password", "", "Password for qBittorrent Web Client")
	qBittorrentSaveFolder := flag.String("qbt-save-folder", "", "Folder to save downloaded torrents")
	qBittorrentCategory := flag.String("qbt-category", "altor-bot", "Set custom category for torrent")

	flag.Parse()

	return &Flags{
		Debug:                      *debug,
		LogLevel:                   *logLevel,
		AnilibriaWSScheme:          *anilibriaWSScheme,
		AnilibriaWSHost:            *anilibriaWSHost,
		AnilibriaWSPath:            *anilibriaWSPath,
		QBittorrentScheme:          *qBittorrentScheme,
		QBittorrentHost:            *qBittorrentHost,
		QBittorrentUser:            *qBittorrentUser,
		QBittorrentPassword:        *qBittorrentPassword,
		QBittorrentSaveFolder:      *qBittorrentSaveFolder,
		QBittorrentIgnoreTSLVerify: *qBittorrentIgnoreTSLVerify,
		QBittorrentCategory:        *qBittorrentCategory,
	}
}

// Config includes all project configurations
type Config struct {
	Debug          bool
	LogLevel       string
	AnilibriaWSURL *AnilibriaWSConfig
	QBittorrent    *QBittorrent
}

// NewConfig Config structure constructor
func NewConfig() (cfg *Config) {
	flags := NewFlags()

	// Switching the logger level to debug mode
	// It is logical that during debugging all logs are needed
	if flags.Debug {
		flags.LogLevel = "debug"
	}

	anilibriaConfig := NewAnilibriaWSConfig(flags)
	qBittorrentConfig := NewQBittorrentConfig(flags)

	return &Config{
		Debug:          flags.Debug,
		LogLevel:       flags.LogLevel,
		AnilibriaWSURL: anilibriaConfig,
		QBittorrent:    qBittorrentConfig,
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
