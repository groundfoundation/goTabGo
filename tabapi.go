package gotabgo

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"strings"

	"github.com/groundfoundation/gotabgo/model"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.WithFields(
		log.Fields{
			"package": "gotabgo",
			"file":    "tabapi.go",
		},
	)
}

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
	credentials := model.Credentials{
		Name:     username,
		Password: password,
		Site: &model.SiteType{
			ContentUrl: contentUrl,
		},
	}

	if impersonateUser != "" {
		credentials.Impersonate = &model.User{
			Name: impersonateUser,
		}
	}
	var tsr model.TsRequest
	tsr.Credentials = credentials

	var payload []byte
	payload, err = getPayload(tsr, t.ContentType)
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
	var tr model.TsResponse
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

func (t *TabApi) NewTrustedTicket(user string, site string) (st string, err error) {
	url := fmt.Sprintf("%s/trusted", t.getUrl())
	//payload := strings.NewReader("username=bjoh121&target_site=RQNS")
	userString := fmt.Sprintf("username=%s&target_site=%s", user, site)
	log.WithField("method", "NewTrustedTicket").
		Debug("userString: ", userString)
	payload := strings.NewReader(userString)
	r, err := http.NewRequest("POST", url, payload)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, _ := http.DefaultClient.Do(r)
	//resp, err := t.c.Post(url, t.ContentType.String(), bytes.NewBuffer(payload))
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	st = string(body)
	log.WithField("method", "NewTrustedTicket").
		Debug("response: ", st)
	return st, nil
}

func (t *TabApi) ServerInfo() (si *model.ServerInfo, err error) {
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
	var tResponse model.TsResponse
	switch t.ContentType {
	case Xml:
		err = xml.Unmarshal(body, &tResponse)
	case Json:
		err = json.Unmarshal(body, &tResponse)
	}
	if err != nil {
		return
	}
	log.WithField("method", "ServerInfo").
		Debug("ServerInfoResponse:\n", tResponse)

	si = &tResponse.ServerInfo

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

// getPayload is a utility function to convert a Go struct into a serialized
// form for a HTTP POST.
func getPayload(thingToEncode interface{}, contentType ContentType) (payload []byte, err error) {
	switch contentType {
	case Xml:
		payload, err = xml.Marshal(thingToEncode)
	case Json:
		payload, err = json.Marshal(thingToEncode)
	}
	return
}

func (t *TabApi) CreateSite(site model.SiteType) (st *model.SiteType, err error) {
	url := fmt.Sprintf("%s/api/%s/sites", t.getUrl(), t.ApiVersion)
	log.WithField("method", "CreateSite").Debug("url: ", string(url))
	var tsRequest model.TsRequest
	tsRequest.Site = site

	var payload []byte
	payload, err = getPayload(tsRequest, t.c.acceptType)
	log.WithField("method", "CreateSite").Debug("payload", string(payload))
	log.WithField("method", "CreateSite").WithField("id", "Token").Debug(t.c.authToken)
	r, e := t.c.Post(url, t.ContentType.String(), bytes.NewBuffer(payload))

	if e != nil {
		log.Error(e)
		return nil, e
	}

	defer r.Body.Close()
	var tResponse model.TsResponse

	body, e := ioutil.ReadAll(r.Body)
	log.WithField("method", "CreateSite").Debug("response", string(body))

	mediaType, _, e := mime.ParseMediaType(r.Header.Get("Content-Type"))
	switch mediaType {
	case "application/xml":
		xml.Unmarshal(body, &tResponse)
	case "application/json":
		json.Unmarshal(body, &tResponse)
	}
	log.WithField("method", "CreateSite").Debug("unmarshal", tResponse)
	if r.StatusCode != http.StatusCreated {
		return nil, &ApiError{r.StatusCode, r.Status}
	}
	log.WithField("method", "CreateSite").Debugf("Error: Code = %d Status = %s", r.StatusCode, r.Status)

	return &tResponse.Site, err
}
