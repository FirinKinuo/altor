package settings

import (
	"net/url"
	"path/filepath"
)

// QBittorrent structure containing QBittorrent web-client configuration
type QBittorrent struct {
	URL        string
	User       string
	Password   string
	SaveFolder string
}

func NewQBittorrentConfig(flags *Flags) (cfg *QBittorrent) {
	qBittorrentURL := url.URL{Scheme: flags.QBittorrentScheme, Host: flags.QBittorrentHost}

	qBittorrentConfig := QBittorrent{
		URL:        qBittorrentURL.String(),
		User:       flags.QBittorrentUser,
		Password:   flags.QBittorrentPassword,
		SaveFolder: filepath.FromSlash(flags.QBittorrentSaveFolder),
	}

	return &qBittorrentConfig
}
