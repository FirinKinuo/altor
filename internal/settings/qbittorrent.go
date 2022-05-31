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

func NewQBittorrentConfig(flags *Flags) (cfg *QBittorrent) {
	qBittorrentURL := url.URL{Scheme: flags.QBittorrentScheme, Host: flags.QBittorrentHost}

	qBittorrentConfig := QBittorrent{
		URL:             qBittorrentURL,
		IgnoreTSLVerify: flags.QBittorrentIgnoreTSLVerify,
		User:            flags.QBittorrentUser,
		Password:        flags.QBittorrentPassword,
		SaveFolder:      filepath.FromSlash(flags.QBittorrentSaveFolder),
		Category:        flags.QBittorrentCategory,
	}

	return &qBittorrentConfig
}
