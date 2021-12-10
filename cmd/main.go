package main

import (
	"auction_auth_service/server"
	"log"
)


func main(){
	app := server.NewApp()

	err := app.Run()
	if err != nil{
		log.Fatal(err)
	}
}