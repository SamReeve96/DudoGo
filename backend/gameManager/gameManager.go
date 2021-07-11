package gameManager

import (
	"fmt"
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

func NewGame(details dudo.NewGameDetails) {

	if details.Players == 0 {
		fmt.Printf("No players, ending \n")
		// No players means the game cant run
		return
	}

	// create a new gamestate
	newGame := game{}
	newGameStatePointer := dudo.SetupGame(details)
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
			for i, game := range activeGames {
				// TO FIX - This is the same game state regardless of the game in the for loop?
				gameState := *game.state
				players := activeGames[i].state.Players
				fmt.Printf("Game ID: %s, Player Count: %v IPlayerCount %v \n", game.id, len(gameState.Players), len(players))
			}
		}
	}
}
