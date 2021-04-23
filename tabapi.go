package gotabgo

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
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
	ctStr, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return err
	}

	contentType, err := ContentTypeString(ctStr)
	if err != nil {
		return err
	}

	var tr model.TsResponse
	err = putResponse(resp.Body, &tr, contentType)
	log.WithField("method", "Signin").
		WithField("id", "unmarshal tr").Debug(tr)
	t.c.authToken = tr.Credentials.Token
	log.WithField("method", "Signin").
		WithField("id", "Token").Debug(t.c.authToken)

	return nil
}

func (t *TabApi) NewTrustedTicket(ttr model.TrustedTicketRequest) (tt model.TrustedTicket, err error) {
	purl := fmt.Sprintf("%s/trusted", t.getUrl())
	data := url.Values{}
	data.Set("username", ttr.Username)
	data.Set("target_site", ttr.Targetsite)
	payload := strings.NewReader(data.Encode())
	var ctype ContentType = Form
	resp, err := t.c.Post(purl, ctype.String(), payload)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = errors.New("Failed: " + resp.Status)
		return
	}
	buf := new(bytes.Buffer)
	io.Copy(buf, resp.Body)
	tt.Value = buf.String()
	log.WithField("method", "NewTrustedTicket").
		Debug("response: ", tt.Value)
	return
}

func (t *TabApi) ServerInfo() (si *model.ServerInfo, err error) {
	//TODO: figure out how to use the apiversion instead of hard coding
	url := fmt.Sprintf("%s/api/%s/serverinfo", t.getUrl(), "2.4")
	r, e := t.c.Get(url)
	if e != nil {
		log.Error(e)
		return nil, e
	}

	log.WithField("method", "ServerInfo").Debug("response:\n", r)
	defer r.Body.Close()
	ctStr, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}
	contentType, err := ContentTypeString(ctStr)
	var tResponse model.TsResponse
	err = putResponse(r.Body, &tResponse, contentType)
	if err != nil {
		return nil, err
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

func putResponse(r io.ReadCloser, dest interface{}, contentType ContentType) (err error) {
	switch contentType {
	case Xml:
		err = xml.NewDecoder(r).Decode(dest)
	case Json:
		err = json.NewDecoder(r).Decode(dest)
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
	r, e := t.Post(url, t.ContentType.String(), bytes.NewBuffer(payload))

	if e != nil {
		log.Error(e)
		return nil, e
	}

	defer r.Body.Close()
	var tResponse model.TsResponse

	mediaType, _, e := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if e != nil {
		return
	}
	contentType, e := ContentTypeString(mediaType)
	if e != nil {
		return
	}
	e = putResponse(r.Body, &tResponse, contentType)
	if e != nil {
		return
	}
	log.WithField("method", "CreateSite").Debug("unmarshal", tResponse)
	if r.StatusCode != http.StatusCreated {
		return nil, &ApiError{r.StatusCode, r.Status}
	}
	log.WithField("method", "CreateSite").Debugf("Error: Code = %d Status = %s", r.StatusCode, r.Status)

	return &tResponse.Site, err
}
