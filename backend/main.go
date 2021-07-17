package main

import (
	"fmt"
	"os"

	"github.com/SamReeve96/DudoGo/backend/cli"
	"github.com/SamReeve96/DudoGo/backend/gameManager"
	"github.com/SamReeve96/DudoGo/backend/server"
)

func main() {
	var option string

	if len(os.Args[1:]) > 0 {
		option = os.Args[1]
	}

	if option == "" {
		option = cli.HandleInput("Run server or build site or both? (1,2,3) fianlly, create a new game? (4) blam")
	}

	switch option {
	case "1":
		serverManager()
	case "2":
		server.BuildApp()
	case "3":
		server.BuildApp()
		serverManager()
	case "4":
		// gameManager.NewGame()
	default:
		fmt.Printf("invalid option")
	}
}

func serverManager() {
	go gameManager.ReportActiveGames()
	server.Serve()
}
