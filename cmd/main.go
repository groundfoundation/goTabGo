package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"

	"github.com/groundfoundation/gotabgo"
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
	}
	log.Debug("Server: ", server)

	user = os.Getenv("TAB_USER")
	if user == "" {
		fmt.Print("Enter User: ")
		user, _ = reader.ReadString('\n')
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
	}
	log.Debug("Password value obtained")

	apiVer = os.Getenv("TABLEAU_API_VERSION")
	if apiVer == "" {
		apiVer = "3.9"
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
	fmt.Printf("Server Info:\n%s", siFmt)
	fmt.Printf("Server Info:\n%s", siXml)

	// Let's login!
	e = tabApi.Signin(user, password, "", "")
	if e != nil {
		log.Fatal(e)
	}
	_, err := tabApi.CreateSite("Test Site 3", "ts3")
	if err != nil {
		log.Error(err.Error())
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
