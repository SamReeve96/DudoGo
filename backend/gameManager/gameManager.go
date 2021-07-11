package gameManager

import (
	"fmt"
	"os"
	"time"

	"github.com/SamReeve96/DudoGo/backend/dudo"

	"github.com/google/uuid"
)

// game holds the pointer to a game state and it's ID
type game struct {
	state *dudo.GameState
	id    uuid.UUID
}

var activeGames []game = []game{}

func NewGame() {
	newGame := game{}
	// create a new gamestate
	newGameStatePointer := dudo.SetupGame()

	if newGameStatePointer == nil {
		fmt.Printf("No players, ending \n")
		// No players means the game cant run
		os.Exit(1)
	}

	newGame.id = uuid.New()
	newGame.state = newGameStatePointer

	fmt.Printf("Main: New game created with ID: %v \n", newGame.id)
	// add to list of games

	activeGames = append(activeGames, newGame)
	fmt.Printf("Main: Currently %v games are running \n", len(activeGames))
}

// While the server is running, report the state of the server to logs
func ReportActiveGames() {
	ticker := time.NewTicker(10 * time.Second)
	for true {
		select {
		case <-ticker.C:
			fmt.Printf("There are currently: %v games active \n", len(activeGames))
			for _, game := range activeGames {
				gameState := *game.state
				fmt.Printf("Game ID: %s, Player Count: %v \n", game.id, len(gameState.Players))
			}
		}
	}
}
