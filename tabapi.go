package gotabgo

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	t.c.Post(url, t.ContentType.String(), bytes.NewBuffer(payload))

	return nil
}

func (t *TabApi) getUrl() string {
	url := "http"
	if t.UseTLS {
		url += "s"
	}
	url += "://" + t.Server

	return url
}
