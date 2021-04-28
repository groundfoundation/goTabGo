package gotabgo

import "errors"

type ContentType int

const (
	Json ContentType = iota
	Xml
	Form
)

func (r ContentType) String() string {
	return [...]string{"application/json", "application/xml", "application/x-www-form-urlencoded"}[r]
}

func ContentTypeString(s string) (c ContentType, e error) {
	var cTypeMap = map[string]ContentType{"application/json": Json, "application/xml": Xml, "application/x-www-form-urlencoded": Form}
	var ok bool

	if c, ok = cTypeMap[s]; ok {
		return
	}
	e = errors.New("Content Type can't be converted")
	return
}

type TabApi struct {
	UseTLS      bool
	Server      string
	ApiVersion  string
	ContentType ContentType
	c           *httpClient
}
