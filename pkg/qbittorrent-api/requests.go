package qbittorrent_api

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

// CreateURLToEndpoint returns full URL up to endpoint
func (c *Client) CreateURLToEndpoint(endpoint string) (requestURL string) {
	URL := url.URL{
		Scheme: c.BaseURL.Scheme,
		Host:   c.BaseURL.Host,
		Path:   fmt.Sprintf("%s/%s", c.BaseURL.Path, endpoint),
	}

	return URL.String()
}

// Post Sends POST request to the specified endpoint
func (c *Client) Post(endpoint string, body *bytes.Buffer, writer *multipart.Writer) (response *http.Response, err error) {
	request, err := http.NewRequest(http.MethodPost, c.CreateURLToEndpoint(endpoint), body)
	if err != nil {
		return nil, fmt.Errorf("error when creating request, reason: %w", err)
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())

	response, err = c.HTTP.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error when sending request, reason: %w", err)
	}

	return response, nil
}

// ReadBodyFromResponse returns the response body as a string
func ReadBodyFromResponse(response *http.Response) (body string, err error) {
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error when reading response, reason: %w", err)
	}

	return string(bodyBytes), nil
}

// AddCookiesToJar If there are more than 0 passed cookies, it will add them to the jar by domain
func (c *Client) AddCookiesToJar(cookies []*http.Cookie) {
	if len(cookies) > 0 {
		c.CookiesJar.SetCookies(&c.BaseURL, cookies)
	}
}

// parseState will return True if ok. is returned in request body
func parseState(body string) (state bool) {
	return strings.Contains(strings.ToLower(body), "ok.")
}
