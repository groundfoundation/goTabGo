package gotabgo

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func NewTabApi(server, version string, useTLS bool, cType ContentType) (*TabApi, error) {
	c := &httpClient{
		acceptType: cType,
	}

	return &TabApi{
		UseTLS:      useTLS,
		Server:      server,
		ApiVersion:  version,
		ContentType: cType,
		c:           c,
	}, nil

}

// Signin authenticates a user and retrieves an auth token
func (t *TabApi) Signin(username, password, contentUrl, impersonateUser string) (err error) {
	url := fmt.Sprintf("%s/api/%s/auth/signin", t.getUrl(), t.ApiVersion)
	credentials := Credentials{
		Name:     username,
		Password: password,
		Site: &Site{
			ContentUrl: contentUrl,
		},
	}

	if impersonateUser != "" {
		credentials.Impersonate = &User{
			Name: impersonateUser,
		}
	}
	signInRequest := SigninRequest{
		Request: credentials,
	}
	var payload []byte
	switch t.ContentType {
	case Xml:
		payload, err = signInRequest.XML()
	case Json:
		payload, err = json.Marshal(signInRequest)
	}
	if err != nil {
		return err
	}
	// Post this to the endpoint
	resp, err := t.c.Post(url, t.ContentType.String(), bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	log.WithField("method", "Signin").Debug(resp)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.WithField("method", "Signin").WithField("id", "body").Debug(string(body))
	var tr TsResponse
	log.Debug("header", resp.Header.Get("Content-Type"))
	contentType, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	switch contentType {
	case "application/xml":
		err = xml.Unmarshal(body, &tr)
	case "application/json":
		err = json.Unmarshal(body, &tr)
	}

	log.WithField("method", "Signin").
		WithField("id", "unmarshal tr").Debug(tr)
	t.c.authToken = tr.Credentials.Token
	log.WithField("method", "Signin").
		WithField("id", "Token").Debug(t.c.authToken)

	return nil
}

func (t *TabApi) ServerInfo() (si *ServerInfo, err error) {
	//TODO: figure out how to use the apiversion instead of hard coding
	url := fmt.Sprintf("%s/api/%s/serverinfo", t.getUrl(), "2.4")
	r, e := t.c.Get(url)
	if e != nil {
		log.Error(e)
		return nil, e
	}

	log.WithField("method", "ServerInfo").
		Debug("response:\n", r)
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	log.WithField("method", "ServerInfo").
		Debug("response:\n", string(body))
	// unmarshal this
	var sir TsResponse
	switch t.ContentType {
	case Xml:
		err = xml.Unmarshal(body, &sir)
	case Json:
		err = json.Unmarshal(body, &sir)
	}
	if err != nil {
		return
	}
	log.WithField("method", "ServerInfo").
		Debug("ServerInfoResponse:\n", sir)

	si = &sir.ServerInfo

	return

}

func (t *TabApi) getUrl() string {
	url := "http"
	if t.UseTLS {
		url += "s"
	}
	url += "://" + t.Server

	return url
}

func (t *TabApi) CreateSite(siteName string) (*Site, error) {
	url := fmt.Sprintf("%s/api/%s/sites", t.getUrl(), t.ApiVersion)
	log.WithField("method", "CreateSite").Debug("url: ", string(url))
	site := Site{
		Name:       siteName,
		ContentUrl: siteName,
	}
	createSiteRequest := CreateSiteRequest{Request: site}
	xmlRep, err := createSiteRequest.XML()
	log.WithField("method", "CreateSite").Debug("xml", xmlRep)
	if err != nil {
		return nil, err
	}
	log.WithField("method", "CreateSite").
		WithField("id", "Token").Debug(t.c.authToken)
	r, e := t.c.Post(url, t.ContentType.String(), bytes.NewBuffer(xmlRep))

	if e != nil {
		log.Error(e)
		return nil, e
	}
	if r.StatusCode != http.StatusCreated {
		return nil, &ApiError{r.StatusCode, r.Status}
	}
	log.WithField("method", "CreateSite").Debugf("Error: Code = %d Status = %s", r.StatusCode, r.Status)
	defer r.Body.Close()
	createSiteResponse := CreateSiteResponse{}
	body, e := ioutil.ReadAll(r.Body)
	log.WithField("method", "CreateSite").Debug("response", string(body))

	return nil, nil
	return &createSiteResponse.Site, err
}
