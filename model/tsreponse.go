package model

import "encoding/xml"

// TsResponse is the wrapper that Tableau Server wraps each response with
type TsResponse struct {
	XMLName     xml.Name    `json:"-" xml:"tsResponse"`
	ServerInfo  ServerInfo  `json:"serverInfo" xml:"serverInfo"`
	Credentials Credentials `json:"credentials" xml:"credentials"`
	Error       ErrorType   `json:"error" xml:"error"`
}

// ServerInfo contains information about product version and api version for the server
type ServerInfo struct {
	XMLName        xml.Name       `json:"-" xml:"serverInfo"`
	ProductVersion ProductVersion `json:"productVersion,omitempty" xml:"productVersion"`
	RestApiVersion string         `json:"restApiVersion,omitempty" xml:"restApiVersion,omitempty"`
}

type SiteType struct {
	XMLName      xml.Name `json:"-" xml:"site"`
	ID           string   `json:"id,omitempty" xml:"id,attr,omitempty"`
	Name         string   `json:"name,omitempty" xml:"name,attr,omitempty"`
	ContentUrl   string   `json:"contentUrl,omitempty" xml:"contentUrl,attr,omitempty"`
	AdminMode    string   `json:"adminMode,omitempty" xml:"adminMode,attr,omitempty"`
	UserQuota    string   `json:"userQuota,omitempty" xml:"userQuota,attr,omitempty"`
	StorageQuota int      `json:"storageQuota,omitempty" xml:"storageQuota,attr,omitempty"`
	State        string   `json:"state,omitempty" xml:"state,attr,omitempty"`
	StatusReason string   `json:"statusReason,omitempty" xml:"statusReason,attr,omitempty"`
	// Usage        *Usage   `json:"usage,omitempty" xml:"usage,omitempty"`
	Usage struct {
		NumUsers     uint `json:"numUsers" xml:"numUsers,attr"`
		NumCreators  uint `json:"numCreators,omitempty" xml:"numCreators,omitempty,attr"`
		NumExplorers uint `json:"numExplorers,omitempty" xml:"numExplorers,omitempty,attr"`
		NumViewers   uint `json:"numViewers,omitempty" xml:"numViewers,omitempty,attr"`
		Storage      uint `json:"storage" xml:"storage,attr"`
	} `json:"usage,omitempty" xml:"usage,omitempty"`
}

type Usage struct {
	XMLName       xml.Name `json:"-" xml:"usage"`
	NumberOfUsers int      `json:"number-of-users" xml:"number-of-users,attr"`
	Storage       int      `json:"storage" xml:"storage,attr"`
}

type Credentials struct {
	XMLName     xml.Name  `json:"-" xml:"credentials"`
	Name        string    `json:"name,omitempty" xml:"name,attr,omitempty"`
	Password    string    `json:"password,omitempty" xml:"password,attr,omitempty"`
	Token       string    `json:"token,omitempty" xml:"token,attr,omitempty"`
	Site        *SiteType `json:"site,omitempty" xml:"site,omitempty"`
	Impersonate *User     `json:"user,omitempty" xml:"user,omitempty"`
}
type ProductVersion struct {
	XMLName xml.Name `json:"-" xml:"productVersion"`
	Value   string   `json:"value" xml:",chardata"`
	Build   string   `json:"build" xml:"build,attr"`
}

type User struct {
	XMLName  xml.Name `json:"-" xml:"user"`
	ID       string   `json:"id,omitempty" xml:"id,attr,omitempty"`
	Name     string   `json:"name,omitempty" xml:"name,attr,omitempty"`
	SiteRole string   `json:"siteRole,omitempty" xml:"siteRole,attr,omitempty"`
	FullName string   `json:"fullName,omitempty" xml:"fullName,attr,omitempty"`
}

type ErrorType struct {
	XMLName xml.Name `json:"-" xml:"error"`
	Summary string   `json:"summary" xml:"summary"`
	Detail  string   `json:"detail" xml:"detail"`
	Code    uint     `json:"code" xml:"code,attr"`
}
