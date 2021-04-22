package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/groundfoundation/gotabgo"
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

	user = os.Getenv("TAB_USER")
	if user == "" {
		fmt.Print("Enter User: ")
		user, _ = reader.ReadString('\n')
	}

	password = os.Getenv("TAB_PASS")
	if password == "" {
		fmt.Print("Enter Password: ")
		pwd, e := term.ReadPassword(int(os.Stdin.Fd()))
		if e != nil {
			log.Fatal(e)
		}
		password = string(pwd)
	}

	fmt.Printf("Server is: %s", server)

	_, e := gotabgo.NewTabApi(server, "2.8", false, gotabgo.Json)
	if e != nil {
		log.Fatal(e)
	}

}
