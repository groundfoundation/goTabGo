package gotabgo

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/groundfoundation/gotabgo/model"
)

const DEFAULT_API_VERSION = "2.5"

type API struct {
	Server    string
	Version   string
	AuthToken string
	SiteName  string
	TlsCert   *TlsCertificate
}

type TlsCertificate struct {
	PrivateKey []byte
	PublicKey  []byte
	Ca         []byte
}

func NewAPI(server string, version string, siteName string, tlsCert *TlsCertificate) *API {
	formatedServername := server
	if strings.HasSuffix(server, "/") {
		formatedServername = server[0 : len(server)-1]
	}
	return &API{Server: formatedServername, Version: version, SiteName: siteName, TlsCert: tlsCert}
}

type SigninRequest struct {
	Request model.Credentials `json:"credentials,omitempty" xml:"credentials,omitempty"`
}

func (req SigninRequest) XML() ([]byte, error) {
	tmp := struct {
		SigninRequest
		XMLName struct{} `xml:"tsRequest"`
	}{SigninRequest: req}
	return xml.MarshalIndent(tmp, "", "   ")
}

type AuthResponse struct {
	Credentials *model.Credentials `json:"credentials,omitempty" xml:"credentials,omitempty"`
}

type ServerInfoResponse struct {
	ServerInfo model.ServerInfo `json:"serverInfo,omitempty" xml:"serverInfo,omitempty"`
}

type CreateSiteRequest struct {
	Request model.SiteType `json:"site,omitempty" xml:"site,omitempty"`
}

func (req CreateSiteRequest) XML() ([]byte, error) {
	tmp := struct {
		CreateSiteRequest
		XMLName struct{} `xml:"tsRequest"`
	}{CreateSiteRequest: req}
	return xml.MarshalIndent(tmp, "", "   ")
}

type CreateSiteResponse struct {
	Site model.SiteType `json:"site,omitempty" xml:"site,omitempty"`
}

type ConnectionCredentials struct {
	Name     string `json:"name,omitempty" xml:"name,attr,omitempty"`
	Password string `json:"password,omitempty" xml:"password,attr,omitempty"`
	Embed    bool   `json:"embed" xml:"embed,attr"`
}

type GoDie struct {
	Code    string `json:"code,omitempty" xml:"code,attr,omitempty"`
	Summary string `json:"summary,omitempty" xml:"summary,omitempty"`
	Detail  string `json:"detail,omitempty" xml:"detail,omitempty"`
}

type ErrorResponse struct {
	Error GoDie `json:"error,omitempty" xml:"error,omitempty"`
}

func (gd GoDie) Error() string {
	return fmt.Sprintf("Code:%s, Summary:%s, Detail:%s", gd.Code, gd.Summary, gd.Detail)
}
