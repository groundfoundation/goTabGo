package gotabgo

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const content_type_header = "Content-Type"
const content_length_header = "Content-Length"
const auth_header = "X-Tableau-Auth"
const application_xml_content_type = "application/json"
const POST = "POST"
const GET = "GET"
const DELETE = "DELETE"

var ErrDoesNotExist = errors.New("Does Not Exist")

func (api *API) Signin(username, password string, contentUrl string, userIdToImpersonate string) error {
	url := fmt.Sprintf("%s/api/%s/auth/signin", api.Server, api.Version)
	credentials := Credentials{Name: username, Password: password}
	if len(userIdToImpersonate) > 0 {
		credentials.Impersonate = &User{ID: userIdToImpersonate}
	}
	siteName := contentUrl
	credentials.Site = &Site{ContentUrl: siteName}
	request := SigninRequest{Request: credentials}
	signInXML, err := request.XML()
	if err != nil {
		return err
	}
	payload := string(signInXML)
	headers := make(map[string]string)
	headers[content_type_header] = application_xml_content_type
	retval := AuthResponse{}
	err = api.makeRequest(url, POST, []byte(payload), &retval, headers, connectTimeOut, readWriteTimeout)
	if err == nil {
		api.AuthToken = retval.Credentials.Token
	}
	return err
}

//http://onlinehelp.tableau.com/current/api/rest_api/en-us/help.htm#REST/rest_api_ref.htm#Sign_Out%3FTocPath%3DAPI%2520Reference%7C_____52
func (api *API) Signout() error {
	url := fmt.Sprintf("%s/api/%s/auth/signout", api.Server, api.Version)
	headers := make(map[string]string)
	headers[content_type_header] = application_xml_content_type
	err := api.makeRequest(url, POST, nil, nil, headers, connectTimeOut, readWriteTimeout)
	return err
}

//http://onlinehelp.tableau.com/current/api/rest_api/en-us/help.htm#REST/rest_api_ref.htm#Server_Info%3FTocPath%3DAPI%2520Reference%7C__
func (api *API) ServerInfo() (ServerInfo, error) {
	// this call only works on apiVersion 2.4 and up
	url := fmt.Sprintf("%s/api/%s/serverinfo", api.Server, "2.4")
	headers := make(map[string]string)
	retval := ServerInfoResponse{}
	err := api.makeRequest(url, GET, nil, &retval, headers, connectTimeOut, readWriteTimeout)
	return retval.ServerInfo, err
}

func (api *API) makeRequest(requestUrl string, method string, payload []byte, result interface{}, headers map[string]string,
	cTimeout time.Duration, rwTimeout time.Duration) error {
	var debug = false
	if debug {
		fmt.Printf("%s:%v\n", method, requestUrl)
		if payload != nil {
			fmt.Printf("%v\n", string(payload))
		}
	}
	client := DefaultTimeoutClient()
	var req *http.Request
	if len(payload) > 0 {
		var httpErr error
		req, httpErr = http.NewRequest(strings.TrimSpace(method), strings.TrimSpace(requestUrl), bytes.NewBuffer(payload))
		if httpErr != nil {
			return httpErr
		}
		req.Header.Add(content_length_header, strconv.Itoa(len(payload)))
	} else {
		var httpErr error
		req, httpErr = http.NewRequest(strings.TrimSpace(method), strings.TrimSpace(requestUrl), nil)
		if httpErr != nil {
			return httpErr
		}
	}
	if headers != nil {
		for header, headerValue := range headers {
			req.Header.Add(header, headerValue)
		}
	}
	if len(api.AuthToken) > 0 {
		if debug {
			fmt.Printf("%s:%s\n", auth_header, api.AuthToken)
		}
		req.Header.Add(auth_header, api.AuthToken)
	}
	var httpErr error
	resp, httpErr := client.Do(req)
	if httpErr != nil {
		return httpErr
	}
	defer resp.Body.Close()
	body, readBodyError := ioutil.ReadAll(resp.Body)
	if debug {
		fmt.Printf("t4g Response:%v\n", string(body))
	}
	if readBodyError != nil {
		return readBodyError
	}
	if resp.StatusCode == 404 {
		return ErrDoesNotExist
	}
	if resp.StatusCode >= 300 {
		tErrorResponse := ErrorResponse{}
		err := xml.Unmarshal(body, &tErrorResponse)
		if err != nil {
			return err
		}
		return tErrorResponse.Error
	}
	if result != nil {
		// else unmarshall to the result type specified by caller
		err := xml.Unmarshal(body, &result)
		if err != nil {
			return err
		}
	}
	return nil
}
