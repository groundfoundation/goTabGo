package gotabgo

import "encoding/xml"

type ContentType int

const (
	Json = iota
	Xml
)

func (r ContentType) String() string {
	return [...]string{"application/json", "application/xml"}[r]
}

type TabApi struct {
	UseTLS      bool
	Server      string
	ApiVersion  string
	ContentType ContentType
	c           *httpClient
}

type TsResponse struct {
	XMLName xml.Name `json:"-" xml:"tsResponse"`

	ServerInfo ServerInfo `json:"serverInfo" xml:"serverInfo"`
}
