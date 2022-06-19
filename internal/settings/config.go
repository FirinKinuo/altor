package settings

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

// Envs structure with all project environments vars
type Envs struct {
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

// NewEnvs - constructor of the Envs structure
// Returns a pointer to the created structure with project environment variables
// Doesn't panic if an invalid variable is specified
// Will log the fatal level with an erroneous parameter
func NewEnvs() (envs *Envs) {
	debug := getEnvBool("DEBUG", false)
	logLevel := getEnvString("LOG_LEVEL", "info")
	anilibriaWSScheme := getEnvString("ANILIBRIA_WS_SCHEME", "wss")
	anilibriaWSHost := getEnvString("ANILIBRIA_WS_HOST", "api.anilibria.tv")
	anilibriaWSPath := getEnvString("ANILIBRIA_WS_PATH", "v2/ws/")
	qBittorrentScheme := getEnvString("QBT_SCHEME", "http")
	qBittorrentHost := getEnvString("QBT_HOST", "localhost:8080")
	qBittorrentUser := getEnvString("QBT_USER", "admin")
	qBittorrentIgnoreTSLVerify := getEnvBool("QBT_IGNORE_TLS_VERIFY", false)
	qBittorrentPassword := getEnvString("QBT_PASSWORD", "")
	qBittorrentSaveFolder := getEnvString("QBT_SAVE_FOLDER", "")
	qBittorrentCategory := getEnvString("QBT_CATEGORY", "altor-bot")

	return &Envs{
		Debug:                      debug,
		LogLevel:                   logLevel,
		AnilibriaWSScheme:          anilibriaWSScheme,
		AnilibriaWSHost:            anilibriaWSHost,
		AnilibriaWSPath:            anilibriaWSPath,
		QBittorrentScheme:          qBittorrentScheme,
		QBittorrentHost:            qBittorrentHost,
		QBittorrentUser:            qBittorrentUser,
		QBittorrentPassword:        qBittorrentPassword,
		QBittorrentSaveFolder:      qBittorrentSaveFolder,
		QBittorrentIgnoreTSLVerify: qBittorrentIgnoreTSLVerify,
		QBittorrentCategory:        qBittorrentCategory,
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
// Returns a pointer to the created structure Config
// Takes data from environment variables
// If variables with the wrong type were specified, it will output an error to the Fatal log with variable key
func NewConfig() (cfg *Config) {
	envs := NewEnvs()

	// Switching the logger level to debug mode
	// It is logical that during debugging all logs are needed
	if envs.Debug {
		envs.LogLevel = "debug"
	}

	anilibriaConfig := NewAnilibriaWSConfig(envs)
	qBittorrentConfig := NewQBittorrentConfig(envs)

	return &Config{
		Debug:          envs.Debug,
		LogLevel:       envs.LogLevel,
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

func getEnvString(key string, default_ string) (value string) {
	value, exists := os.LookupEnv(key)

	if exists {
		return value
	}

	return default_
}

func getEnvBool(key string, default_ bool) (value bool) {
	valueString, exist := os.LookupEnv(key)

	if exist {
		value, err := strconv.ParseBool(valueString)
		if err != nil {
			log.Fatalf("env %s is not a bool", key)
		}

		return value
	}

	return default_
}
