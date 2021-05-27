package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"github.com/groundfoundation/gotabgo"
	"github.com/groundfoundation/gotabgo/model"
	log "github.com/sirupsen/logrus"
	"golang.org/x/term"
)

func main() {
	var user, server, password, apiVer string
	reader := bufio.NewReader(os.Stdin)

	server = os.Getenv("TAB_SERVER")
	if server == "" {
		fmt.Print("Enter Server: ")
		server, _ = reader.ReadString('\n')
		server = strings.Replace(server, "\n", "", -1)
	}
	log.Debug("Server: ", server)

	user = os.Getenv("TAB_USER")
	if user == "" {
		fmt.Print("Enter User: ")
		user, _ = reader.ReadString('\n')
		user = strings.Replace(user, "\n", "", -1)
	}
	log.Debug("User: ", user)

	password = os.Getenv("TAB_PASS")
	if password == "" {
		fmt.Print("Enter Password: ")
		pwd, e := term.ReadPassword(int(os.Stdin.Fd()))
		if e != nil {
			log.Fatal(e)
		}
		password = string(pwd)
		password = strings.Replace(password, "\n", "", -1)
	}
	log.Debug("Password value obtained")

	apiVer = os.Getenv("TABLEAU_API_VERSION")
	if apiVer == "" {
		apiVer = "3.7"
	}
	log.Debug("API Version:", apiVer)

	fmt.Printf("\nServer is: %s", server)

	tabApi, e := gotabgo.NewTabApi(server, apiVer, true, gotabgo.Xml)
	if e != nil {
		log.Fatal(e)
	}

	log.Debug("tabApi", tabApi)
	si, e := tabApi.ServerInfo()
	if e != nil {
		log.Fatal(e)
	}
	siFmt, _ := json.MarshalIndent(si, "", "\t")
	siXml, _ := xml.MarshalIndent(si, "", "\t")
	fmt.Printf("\nJSON Server Info:\n%s", siFmt)
	fmt.Printf("\nXMLServer Info:\n%s", siXml)

	// Let's login!
	e = tabApi.Signin(user, password, "", "")
	if e != nil {
		log.Fatal(e)
	}
	// Let's Create a site!
	fmt.Print("\nAbout to test Create Site Fuction.\nEnter Site Name: ")
	siteName, _ := reader.ReadString('\n')
	fmt.Print("Enter Site URL ID: ")
	siteURL, _ := reader.ReadString('\n')
	site, err := tabApi.CreateSite(model.SiteType{Name: siteName, ContentUrl: siteURL})
	if err != nil {
		log.Error(err.Error())
	}
	j, _ := json.Marshal(site)
	fmt.Print(string(j))

	// Let's get a trusted ticket!
	fmt.Print("\nAbout to test trusted ticket.\nEnter User Name for ticket: ")
	userName, _ := reader.ReadString('\n')
	userName = strings.TrimSuffix(userName, "\n")
	fmt.Print("\nAbout to test trusted ticket.\nEnter Site: ")
	tSite, _ := reader.ReadString('\n')
	tSite = strings.TrimSuffix(tSite, "\n")
	ticket, err := tabApi.NewTrustedTicket(model.TrustedTicketRequest{Username: userName, Targetsite: tSite})
	if err != nil {
		log.Error(err.Error())
	}
	fmt.Printf("\nTicket: %s\n", ticket.Value)

	// Let's list reports for a user
	fmt.Printf("\nChecking for user: %s\n", userName)
	userStruct, err := tabApi.QueryUserOnSite(userName)
	if err != nil {
		log.Error(err.Error())
	}
	fmt.Printf("\nListing Reports for ID: %s\n", userStruct.ID)
	workbooks, err := tabApi.ListReportsForUser(userStruct)
	if err != nil {
		log.Error(err.Error())
	}
	for _, workbook := range workbooks {
		fmt.Printf("\nWorkbook: %v\n", workbook.Name)
	}
}

func init() {
	var dev bool
	// Determine where we are running
	if _, exists := os.LookupEnv("DEVELOPMENT"); exists {
		dev = true
	}

	if dev {
		log.SetLevel(log.DebugLevel)
	}
}
