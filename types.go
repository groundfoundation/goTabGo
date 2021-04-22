package gotabgo

type ContentType int

const (
	Xml = iota
	Json
)

func (r ContentType) String() string {
	return [...]string{"application/xml", "application/json"}[r]
}

type TabApi struct {
	UseTLS      bool
	Server      string
	ApiVersion  string
	ContentType ContentType
	c           *httpClient
}
