package gotabgo

import (
	"io"
	"net/http"
)

const (
	TABLEAU_AUTH_HEADER = "X-Tableau-Auth"
)

type httpClient struct {
	c          http.Client
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
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return c.c.Do(req)
}

func (c *httpClient) Do(req *http.Request) (*http.Response, error) {
	if c.authToken != "" {
		req.Header.Add(TABLEAU_AUTH_HEADER, c.authToken)
	}
	req.Header.Add("Accept", c.acceptType.String())
	return c.c.Do(req)
}
