package main

import (
	"github.com/debidarmawan/debozero-backend/api"
)

func main() {
	server := api.NewServer(".")
	server.Start(8000)
}
