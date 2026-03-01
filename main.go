package main

import (
	"go-final-project/pkg/server"
	"log"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}