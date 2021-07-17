package gameManager

import (
	"errors"
	"fmt"
	"time"

	"github.com/SamReeve96/DudoGo/backend/dudo"

	"github.com/google/uuid"
)

// game holds the game state and it's ID
type Game struct {
	state      dudo.GameState
	id         uuid.UUID
	friendlyID string
}

var activeGames []Game = []Game{}

func NewGame(details dudo.NewGameDetails) {

	if details.Players == 0 {
		fmt.Printf("No players, ending \n")
		// No players means the game cant run
		return
	}

	// create a new gamestate
	newGame := Game{}
	newGameState := dudo.SetupGame(details)
	newGame.id = uuid.New()
	newGame.state = newGameState
	newGame.friendlyID = details.FriendlyID

	fmt.Printf("Main: New game created with ID: %v \n", newGame.id)
	// add to list of games

	activeGames = append(activeGames, newGame)
	fmt.Printf("Main: Currently %v games are running \n", len(activeGames))
}

func JoinGame(friendlyID string, playerName string) bool {
	gameToJoin := getGame(friendlyID)

	if gameToJoin.id == uuid.Nil {
		return false
	}

	if tryAddPlayer(gameToJoin, playerName) {
		return true
	}
	// Didnt add player, game is full (all players have names)

	return false
}

func tryAddPlayer(game Game, newPlayerName string) bool {
	if game.state.Started {
		//cant join game in progress
		return false
	}
	for i, player := range game.state.Players {
		if player.Name == newPlayerName {
			// player in game already has name
			return false
		}
		if player.Name == "" {
			game.state.Players[i].Name = newPlayerName
			// TODO: replace lazy bools with errs and handle better
			return true
		}
	}
	return false
}

func StartGame(friendlyID string, playerName string) error {
	gameToStart := getGame(friendlyID)

	if gameToStart.state.Players[1].Name == "" {
		return errors.New("there is only the creator present, cant start game")
	}

	if gameToStart.state.Players[0].Name != playerName {
		return errors.New("only the creator can start a game")
	}

	gameToStart.state.Started = true

	go dudo.RunGame()

	return nil
}

func getGame(friendlyID string) Game {
	for _, game := range activeGames {
		if game.friendlyID == friendlyID {
			return game
		}
	}
	return Game{}
}

// While the server is running, report the state of the server to logs
func ReportActiveGames() {
	ticker := time.NewTicker(10 * time.Second)
	for true {
		select {
		case <-ticker.C:
			fmt.Printf("There are currently: %v games active \n", len(activeGames))
			for i, game := range activeGames {
				players := activeGames[i].state.Players
				fmt.Printf("------------------------------------------------ \n")
				fmt.Printf("Game ID: %s \n", game.id)
				fmt.Printf("Player Count: %v \n", len(players))
				fmt.Printf("friendlyID: %s \n", game.friendlyID)
				fmt.Printf("started: %v \n", game.state.Started)
				for _, player := range players {
					fmt.Printf("Player name: %s \n", player.Name)
				}
				fmt.Printf("------------------------------------------------ \n")
			}
		}
	}
}

func FriendlyIDInUse(newID string) bool {
	for _, game := range activeGames {
		if game.friendlyID == newID {
			return true
		}
	}
	return false
}
