package qbittorrent_api

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"strings"
)

// createLoginPayload returns the writer and bytes.Buffer from the given login data
func createLoginPayload(username, password string) (writer *multipart.Writer, payload *bytes.Buffer, err error) {
	payload = &bytes.Buffer{}
	writer = multipart.NewWriter(payload)

	err = writer.WriteField("username", username)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot add field username, reason: %w", err)
	}

	err = writer.WriteField("password", password)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot add field password, reason: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return nil, nil, fmt.Errorf("error when closing writer: %w", err)
	}

	return writer, payload, nil
}

// parseIsLogin return true if authorization was successful
func parseIsLogin(body string) (state bool) {
	return strings.Contains(strings.ToLower(body), "ok.")
}

// Login sends an authorization request in the qBittorrent web client
func (c *Client) Login(username, password string) (isLogin bool, err error) {
	writer, payload, err := createLoginPayload(username, password)
	if err != nil {
		return false, fmt.Errorf("error when creating payload, reason: %w", err)
	}

	response, err := c.Post("auth/login", payload, writer)
	if err != nil {
		return false, fmt.Errorf("error when sending POST query, reason: %w", err)
	}

	body, err := ReadBodyFromResponse(response)
	if err != nil {
		return false, err
	}

	isLogin = parseIsLogin(body)

	if isLogin == false {
		return false, fmt.Errorf("authentification failed, reason: %s", body)
	}

	c.AddCookiesToJar(response.Cookies())
	c.SetAuthorizedStatus(true)

	return isLogin, nil
}
