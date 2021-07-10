package main

import (
	"fmt"
	"os"

	"github.com/SamReeve96/DudoGo/logic"
)

// serverState - holds states of games
type serverState struct {
	activeGames []*logic.GameState
}

func main() {
	server := serverState{
		// change to pointers so that we can access live state rather than instance of object?
		activeGames: []*logic.GameState{},
	}

	fmt.Printf("Main: Starting game logic \n")

	// create a new gamestate
	newGamePointer := logic.SetupGame()
	var newGame logic.GameState

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
	gameToRun := *server.activeGames[0]
	gameToRun.RunGame()
}
