package model

import "encoding/xml"

// TsRequest is the wrapper that Tableau Server expects requests to be wrapped with
type TsRequest struct {
	XMLName     xml.Name    `json:"-" xml:"tsRequest"`
	Credentials Credentials `json:"credentials,omitempty" xml:"credentials,omitempty"`
}
