package main

import (
	"fmt"
	"os"

	"github.com/Zhousiru/inker/api"
	"github.com/Zhousiru/inker/db"
)

func setup() {
	var username string
	var password string

	for {
		fmt.Println("Please enter username:")
		fmt.Scanln(&username)
		if username != "" {
			break
		}
		fmt.Println("Illegal username.")
	}

	for {
		fmt.Println("Please enter password:")
		fmt.Scanln(&password)
		if password != "" {
			break
		}
		fmt.Println("Illegal password.")
	}

	err := db.NewUser(username, password)
	if err != nil {
		fmt.Println("Unable to create user: ", err)
		return
	}
	fmt.Println("User created.")
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "setup" {
		setup()
		return
	}
	api.Init()
}
