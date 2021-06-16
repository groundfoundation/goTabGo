package gotabgo

import (
	"io"
	"net"
	"net/http"
	"strings"

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

func (c *httpClient) PostWithIP(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	localIp := GetOutboundIP(url).String()
	log.WithField("Method", "httpclient.PostWithIP").Debugf("LocalIP: %s", localIp)
	log.Printf("LocalIP: %s", localIp)
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-FORWARDED-FOR", localIp)
	req.Header.Set("Content-Type", contentType)
	return c.Do(req)
}

func (c *httpClient) Do(req *http.Request) (*http.Response, error) {
	if c.authToken != "" {
		req.Header.Add(TABLEAU_AUTH_HEADER, c.authToken)
	}
	req.Header.Add("Accept", c.acceptType.String())
	return c.client.Do(req)
}

func GetOutboundIP(dialAddress string) net.IP {
	log.WithField("Method", "httpclient.GetOutboundIP").Debugf("outboundURL: %s", dialAddress)
	splitAddy := strings.Split(dialAddress, "/")
	conn, err := net.Dial("udp", splitAddy[2]+":80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
