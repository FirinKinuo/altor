package main

import (
	"altor/internal/settings"
	"altor/internal/websocket"
	"altor/pkg/qbittorrent-api"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var cfg = settings.NewConfig()

func AuthQBittorrentClient(qbtClient *qbittorrent_api.Client) {
	isLogin, err := qbtClient.Login(cfg.QBittorrent.User, cfg.QBittorrent.Password)
	if err != nil {
		log.Panicf("unable to login qbittorrent client, reason: %s", err)
	}

	if isLogin == false {
		log.Panicf("unable to login qbittorrent client via %s, reason: invalid credentials", cfg.QBittorrent.User)
	}
}

// HandlingTorrentsDownloading Creates a handler that waits for data from a torrentsChannel
//
// It's best to run it like a goroutine
func HandlingTorrentsDownloading(qbtClient *qbittorrent_api.Client, torrentsChannel <-chan []string) {
	torrents := qbittorrent_api.DownloadTorrents{
		Urls:       []string{},
		SaveFolder: cfg.QBittorrent.SaveFolder,
		Category:   cfg.QBittorrent.Category,
	}

	for {
		torrents.Urls = <-torrentsChannel // Write an array of links to torrents from the channel torrentsChannel
		log.Infof("Received a request to download a torrents: %s, send it to qBittorrent", torrents.Urls)

		// Use goroutine to unlock the main thread and not wait for a response from qBittorrent
		go func() {
			AuthQBittorrentClient(qbtClient) // Reauthorization to update the session cookie

			status, err := qbtClient.Download(&torrents)
			if err != nil {
				log.Error(err)
				return
			}
			if status != http.StatusOK {
				log.Error(fmt.Sprintf("torrent download request failed, response code %d should be 200", status))
				return
			}

			log.Info("Torrent successfully submitted for download")
		}()
	}
}

func main() {
	settings.ConfigureLogger(cfg.LogLevel)
	log.Info("Starting Altor!")

	qbt, err := qbittorrent_api.NewClient(cfg.QBittorrent.URL, cfg.QBittorrent.IgnoreTSLVerify)
	if err != nil {
		log.Fatalf("unable to create a client to work with qBittorrent, reason: %s", err)
	}

	AuthQBittorrentClient(qbt)
	log.Info("Successful authorization in the qBittorrent client")

	torrentsChannel := make(chan []string)

	go HandlingTorrentsDownloading(qbt, torrentsChannel)

	anilibriaWS, err := websocket.NewAnilibriaWS(cfg.AnilibriaWSURL.GetURL())
	if err != nil {
		log.Panicln(err)
	}

	anilibriaWS.AttachTorrentsChannel(torrentsChannel)

	err = anilibriaWS.ListenSubscribes()
	if err != nil {
		log.Panicln(err)
	}
}
