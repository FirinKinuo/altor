package qbittorrent_api

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
)

type DownloadTorrents struct {
	Urls       []string
	SaveFolder string
	Category   string
}

// renderDownloadUrlsString combines all url in an array into a single string separated by \n
func renderDownloadUrlsString(urlsArray []string) (urls string) {
	return strings.Join(urlsArray, "\n")
}

// createDownloadPayload returns the writer and bytes.Buffer from the given DownloadTorrents
func createDownloadPayload(torrents *DownloadTorrents) (writer *multipart.Writer, payload *bytes.Buffer, err error) {
	payload = &bytes.Buffer{}
	writer = multipart.NewWriter(payload)

	torrentUrls := renderDownloadUrlsString(torrents.Urls)

	err = writer.WriteField("urls", torrentUrls)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot add field urls, reason: %w", err)
	}

	err = writer.WriteField("savepath", torrents.SaveFolder)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot add field savepath, reason: %w", err)
	}

	err = writer.WriteField("category", torrents.Category)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot add field category, reason: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return nil, nil, fmt.Errorf("error when closing writer: %w", err)
	}

	return writer, payload, nil
}

// Download sends a request to download the torrent URLs specified in DownloadTorrents
//
// Returns a status code to be able to trace the issue when trying to download
func (c *Client) Download(torrents *DownloadTorrents) (statusCode int, err error) {
	writer, payload, err := createDownloadPayload(torrents)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("error when creating download payload, reason: %w", err)
	}

	response, err := c.Post("torrents/add", payload, writer)
	if err != nil {
		return http.StatusForbidden, fmt.Errorf("error when sending POST query, reason: %w", err)
	}

	if response.StatusCode == http.StatusForbidden {
		return http.StatusForbidden, fmt.Errorf("request was sent without authorization")
	}

	body, err := ReadBodyFromResponse(response)
	if err != nil {
		return http.StatusForbidden, fmt.Errorf("unable to read body from request, reason: %w", err)
	}

	state := parseState(body)

	// If ok. was not sent to the response body, then something went wrong
	if state == false {
		return http.StatusForbidden, fmt.Errorf("downloading request failed, reason: %s", body)
	}

	return http.StatusOK, nil
}
