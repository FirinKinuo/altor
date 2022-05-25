package settings

import (
	"net/url"
)

// AnilibriaWSConfig содержит URL для подключения к Websocket Anilibria
type AnilibriaWSConfig struct {
	Scheme string
	Host   string
	Path   string
}

func NewAnilibriaWSConfig(flags *Flags) (cfg *AnilibriaWSConfig) {
	anilibriaWSConfig := AnilibriaWSConfig{
		Scheme: flags.AnilibriaWSScheme,
		Host:   flags.AnilibriaWSHost,
		Path:   flags.AnilibriaWSPath,
	}

	return &anilibriaWSConfig
}

func (w *AnilibriaWSConfig) GetURL() string {
	url_ := url.URL{Scheme: w.Scheme, Host: w.Host, Path: w.Path}

	return url_.String()
}

func (w *AnilibriaWSConfig) GetURLWithMethod(method string) string {
	url_ := url.URL{
		Scheme: w.Scheme,
		Host:   w.Host,
		Path:   w.Path + method,
	}

	return url_.String()
}
