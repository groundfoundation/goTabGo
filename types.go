package gotabgo

type ContentType int

const (
	Json ContentType = iota
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
