package gotabgo

type ContentType int

const (
	Xml = iota
	Json
)

func (r ContentType) String() string {
	return [...]string{"application/xml", "application/json"}[r]
}
