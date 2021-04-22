package gotabgo

import (
	"encoding/xml"
	"fmt"
	"strings"
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
	Request Credentials `json:"credentials,omitempty" xml:"credentials,omitempty"`
}

func (req SigninRequest) XML() ([]byte, error) {
	tmp := struct {
		SigninRequest
		XMLName struct{} `xml:"tsRequest"`
	}{SigninRequest: req}
	return xml.MarshalIndent(tmp, "", "   ")
}

type AuthResponse struct {
	Credentials *Credentials `json:"credentials,omitempty" xml:"credentials,omitempty"`
}

type ServerInfoResponse struct {
	ServerInfo ServerInfo `json:"serverInfo,omitempty" xml:"serverInfo,omitempty"`
}

type ProductVersion struct {
	Value string `json:"value"`
	Build string `json:"build"`
}

type ServerInfo struct {
	ProductVersion ProductVersion `json:"productVersion,omitempty" xml:"productVersion,omitempty"`
	RestApiVersion string         `json:"restApiVersion,omitempty" xml:"restApiVersion,omitempty"`
}

type Credentials struct {
	Name        string `json:"name,omitempty" xml:"name,attr,omitempty"`
	Password    string `json:"password,omitempty" xml:"password,attr,omitempty"`
	Token       string `json:"token,omitempty" xml:"token,attr,omitempty"`
	Site        *Site  `json:"site,omitempty" xml:"site,omitempty"`
	Impersonate *User  `json:"user,omitempty" xml:"user,omitempty"`
}

type User struct {
	ID       string `json:"id,omitempty" xml:"id,attr,omitempty"`
	Name     string `json:"name,omitempty" xml:"name,attr,omitempty"`
	SiteRole string `json:"siteRole,omitempty" xml:"siteRole,attr,omitempty"`
	FullName string `json:"fullName,omitempty" xml:"fullName,attr,omitempty"`
}

type Site struct {
	ID           string     `json:"id,omitempty" xml:"id,attr,omitempty"`
	Name         string     `json:"name,omitempty" xml:"name,attr,omitempty"`
	ContentUrl   string     `json:"contentUrl,omitempty" xml:"contentUrl,attr,omitempty"`
	AdminMode    string     `json:"adminMode,omitempty" xml:"adminMode,attr,omitempty"`
	UserQuota    string     `json:"userQuota,omitempty" xml:"userQuota,attr,omitempty"`
	StorageQuota int        `json:"storageQuota,omitempty" xml:"storageQuota,attr,omitempty"`
	State        string     `json:"state,omitempty" xml:"state,attr,omitempty"`
	StatusReason string     `json:"statusReason,omitempty" xml:"statusReason,attr,omitempty"`
	Usage        *SiteUsage `json:"usage,omitempty" xml:"usage,omitempty"`
}

type SiteUsage struct {
	NumberOfUsers int `json:"number-of-users" xml:"number-of-users,attr"`
	Storage       int `json:"storage" xml:"storage,attr"`
}

type CreateSiteRequest struct {
	Request Site `json:"site,omitempty" xml:"site,omitempty"`
}

func (req CreateSiteRequest) XML() ([]byte, error) {
	tmp := struct {
		CreateSiteRequest
		XMLName struct{} `xml:"tsRequest"`
	}{CreateSiteRequest: req}
	return xml.MarshalIndent(tmp, "", "   ")
}

type CreateSiteResponse struct {
	Site Site `json:"site,omitempty" xml:"site,omitempty"`
}

type ConnectionCredentials struct {
	Name     string `json:"name,omitempty" xml:"name,attr,omitempty"`
	Password string `json:"password,omitempty" xml:"password,attr,omitempty"`
	Embed    bool   `json:"embed" xml:"embed,attr"`
}

func (s Site) XML() ([]byte, error) {
	return xml.MarshalIndent(s, "", "   ")
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
