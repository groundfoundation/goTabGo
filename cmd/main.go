package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/groundfoundation/gotabgo"
	log "github.com/sirupsen/logrus"
	"golang.org/x/term"
)

func main() {
	var user, server, password string
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

	fmt.Printf("\nServer is: %s", server)

	tabApi, e := gotabgo.NewTabApi(server, "2.8", true, gotabgo.Json)
	if e != nil {
		log.Fatal(e)
	}

	log.Debug("tabApi", tabApi)
	tabApi.ServerInfo()

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
