package settings

import (
	"net/url"
	"path/filepath"
)

// QBittorrent structure containing QBittorrent web-client configuration
type QBittorrent struct {
	URL             url.URL
	IgnoreTSLVerify bool
	User            string
	Password        string
	SaveFolder      string
	Category        string
}

func NewQBittorrentConfig(envs *Envs) (cfg *QBittorrent) {
	qBittorrentURL := url.URL{Scheme: envs.QBittorrentScheme, Host: envs.QBittorrentHost}

	qBittorrentConfig := QBittorrent{
		URL:             qBittorrentURL,
		IgnoreTSLVerify: envs.QBittorrentIgnoreTSLVerify,
		User:            envs.QBittorrentUser,
		Password:        envs.QBittorrentPassword,
		SaveFolder:      filepath.FromSlash(envs.QBittorrentSaveFolder),
		Category:        envs.QBittorrentCategory,
	}

	return &qBittorrentConfig
}
