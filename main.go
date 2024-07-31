package main

import (
	"github.com/jnaraujo/mcprotocol/server"
)

func main() {
	sv := server.NewServer(":25565")

	err := sv.Listen()
	if err != nil {
		panic(err)
	}
}
