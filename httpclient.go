package gotabgo

import (
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const (
	TABLEAU_AUTH_HEADER = "X-Tableau-Auth"
)

type httpClient struct {
	client     http.Client
	authToken  string
	acceptType ContentType
}

func (c *httpClient) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *httpClient) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return c.Do(req)
}

func (c *httpClient) Do(req *http.Request) (*http.Response, error) {
	log.Debug("Header: ", TABLEAU_AUTH_HEADER, c.authToken)
	if c.authToken != "" {
		req.Header.Add(TABLEAU_AUTH_HEADER, c.authToken)
	}
	req.Header.Add("Accept", c.acceptType.String())
	return c.client.Do(req)
}
