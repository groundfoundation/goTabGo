package model

import "encoding/xml"

// TsRequest is the wrapper that Tableau Server expects requests to be wrapped with
type TsRequest struct {
	XMLName     xml.Name    `json:"-" xml:"tsRequest"`
	Credentials Credentials `json:"credentials,omitempty" xml:"credentials,omitempty"`
	Site        SiteType    `json:"site,omitempty" xml:"site,omitempty"`
}

//
type TrustedTicket struct {
	Username    string `json:"username,omitempty" xml:"username,omitempty"`
	Target_site string `json:"target_site,omitempty" xml:"target_site,omitempty"`
}
