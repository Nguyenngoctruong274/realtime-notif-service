package main

import (
	"log"
	"yes4all/ads-noti-api/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Printf("error while execute: %s", err.Error())
	}
}
