package torrent

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type (
	// Torrent data structure
	// Indicates the type of content and its download link
	Torrent struct {
		ID   uint32
		Type string
		Hash string
	}

	// Message with anime ID in Anilibria database and Torrent list
	Message struct {
		AnimeID     uint16
		TorrentList []Torrent
	}

	// Response message from WebSocket
	Response struct {
		TorrentUpdate struct {
			ID       string `json:"id"`
			Torrents struct {
				List []struct {
					TorrentID int `json:"torrent_id"`
					Series    struct {
						String string `json:"string"`
					} `json:"series"`
					Quality struct {
						String  string `json:"string"`
						Encoder string `json:"encoder"`
					} `json:"quality"`
					TotalSize int64  `json:"total_size"`
					Hash      string `json:"hash"`
				} `json:"list"`
			} `json:"torrents"`
		} `json:"torrent_update"`
	}
)

// GetAnimeID - getter the anime ID from the response
func (r *Response) GetAnimeID() (id uint16, err error) {
	animeID, err := strconv.ParseUint(r.TorrentUpdate.ID, 0, 16)
	if err != nil {
		return 0, fmt.Errorf("unable to get anime ID from response, reason: %w", err)
	}

	return uint16(animeID), nil
}

// GetTorrents - getter qbittorrent-api list from request
func (r *Response) GetTorrents() (torrents []Torrent, err error) {
	var torrentList []Torrent

	for _, torrent := range r.TorrentUpdate.Torrents.List {
		torrentList = append(torrentList, Torrent{
			ID:   uint32(torrent.TorrentID),
			Type: fmt.Sprintf("%s [%s] %s", torrent.Series.String, torrent.Quality.String, torrent.Quality.Encoder),
			Hash: torrent.Hash,
		})
	}
	return torrentList, nil
}

// Parse connection status message
//
// Takes a response in bytes as an argument
func Parse(response []byte) (m *Message, err error) {
	var jsonResponse Response

	err = json.Unmarshal(response, &jsonResponse)
	if err != nil {
		return nil, fmt.Errorf("unable to read json TorrentUpdate, reason: %w", err)
	}

	animeID, err := jsonResponse.GetAnimeID()
	if err != nil {
		log.Error(err)
	}

	torrents, err := jsonResponse.GetTorrents()
	if err != nil {
		return nil, err
	}

	return &Message{
		AnimeID:     animeID,
		TorrentList: torrents,
	}, nil
}

// GetTorrentsHash returns a list of received torrent hashes
func (m *Message) GetTorrentsHash() (torrentsURIList []string) {
	for _, torrent := range m.TorrentList {
		torrentsURIList = append(torrentsURIList, torrent.Hash)
	}

	return torrentsURIList
}
