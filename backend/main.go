package main

import (
	"fmt"
	"os"

	"github.com/SamReeve96/DudoGo/backend/cli"
	"github.com/SamReeve96/DudoGo/backend/dudo"
	"github.com/SamReeve96/DudoGo/backend/server"
)

// serverState - holds states of games
type serverState struct {
	activeGames []*dudo.GameState
}

func main() {
	option := cli.HandleInput("Run server or game or both? (1,2,3)")

	switch option {
	case "1":
		serverManager()
	case "2":
		gameManager()
	case "3":
		fmt.Printf("Comming soon tm")
	default:
		fmt.Printf("invalid option")
	}
}

func serverManager() {
	server.Serve()
}

func gameManager() {
	server := serverState{
		// change to pointers so that we can access live state rather than instance of object?
		activeGames: []*dudo.GameState{},
	}

	fmt.Printf("Main: Starting game logic \n")

	// create a new gamestate
	newGamePointer := dudo.SetupGame()
	var newGame dudo.GameState

	if newGamePointer == nil {
		fmt.Printf("No players, ending \n")
		// No players means the game cant run
		os.Exit(1)
	} else {
		newGame = *newGamePointer
	}

	fmt.Printf("Main: New game created featuring %v players \n", len(newGame.Players))
	// add to list of games
	server.activeGames = append(server.activeGames, newGamePointer)
	fmt.Printf("Main: Currently %v games are running \n", len(server.activeGames))

	// get the game to being running it
	//gameToRun := *server.activeGames[0]

	// this only runs the game being setup. would be better to call run game from game state (create a game interface?)
	dudo.RunGame()
}
