package qbittorrent_api

import (
	"crypto/tls"
	"fmt"
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

const (
	ApiPath string = "api/v2"
)

// Client creates a connection to qbittorrent and performs requests
type Client struct {
	HTTP       *http.Client
	BaseURL    url.URL
	CookiesJar http.CookieJar
	Authorized bool
}

// NewClient - Client structure constructor
func NewClient(url url.URL, ignoreTLS bool) (qbtClient *Client, err error) {
	url.Path = ApiPath

	cookieJar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, fmt.Errorf("unable to create cookiejar, reason: %w", err)
	}

	HTTP := &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: ignoreTLS}},
		Jar:       cookieJar,
	}

	return &Client{
		HTTP:       HTTP,
		BaseURL:    url,
		CookiesJar: cookieJar,
	}, nil
}

func (c *Client) SetAuthorizedStatus(status bool) {
	c.Authorized = status
}
