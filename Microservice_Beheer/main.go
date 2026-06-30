package main

import (
	//"beheer/database"
	"log"
)

func main() {
	//database.Init()

	server := NewAPIServer(":3000")
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
