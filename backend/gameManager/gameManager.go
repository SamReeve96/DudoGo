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
	State      dudo.GameState
	Id         uuid.UUID
	FriendlyID string
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
	newGame.Id = uuid.New()
	newGame.State = newGameState
	newGame.FriendlyID = details.FriendlyID

	fmt.Printf("Main: New game created with ID: %v \n", newGame.Id)
	// add to list of games

	activeGames = append(activeGames, newGame)
	fmt.Printf("Main: Currently %v games are running \n", len(activeGames))
}

func JoinGame(friendlyID string, playerName string) bool {
	gameToJoin := GetGame(friendlyID)

	if gameToJoin.Id == uuid.Nil {
		return false
	}

	if tryAddPlayer(gameToJoin, playerName) {
		return true
	}
	// Didnt add player, game is full (all players have names)

	return false
}

func tryAddPlayer(game Game, newPlayerName string) bool {
	if game.State.Started {
		//cant join game in progress
		return false
	}
	for i, player := range game.State.Players {
		if player.Name == newPlayerName {
			// player in game already has name
			return false
		}
		if player.Name == "" {
			game.State.Players[i].Name = newPlayerName
			// TODO: replace lazy bools with errs and handle better
			return true
		}
	}
	return false
}

func UpdateGame(GameDetails Game) bool {
	for i := range activeGames {
		if GameDetails.FriendlyID == GameDetails.FriendlyID {
			activeGames[i] = GameDetails
			return true
		}
	}
	return false
}

func StartGame(gameToStart dudo.GameState, playerName string) (dudo.GameState, error) {
	if gameToStart.Players[1].Name == "" {
		return gameToStart, errors.New("there is only the creator present, cant start game")
	}

	if gameToStart.Players[0].Name != playerName {
		return gameToStart, errors.New("only the creator can start a game")
	}

	gameToStart.Started = true

	// start the first round
	gameToStart = dudo.NewRound(gameToStart)

	return gameToStart, nil
}

func GetGame(friendlyID string) Game {
	for _, game := range activeGames {
		if game.FriendlyID == friendlyID {
			return game
		}
	}
	return Game{}
}

// While the server is running, report the state of the server to logs
func ReportActiveGames() {
	ticker := time.NewTicker(5 * time.Second)
	for true {
		select {
		case <-ticker.C:
			fmt.Printf("There are currently: %v games active \n", len(activeGames))
			for i, game := range activeGames {
				players := activeGames[i].State.Players
				fmt.Printf("------------------------------------------------ \n")
				fmt.Printf("Game ID: %s \n", game.Id)
				fmt.Printf("Player Count: %v \n", len(players))
				fmt.Printf("friendlyID: %s \n", game.FriendlyID)
				fmt.Printf("started: %v \n", game.State.Started)
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
		if game.FriendlyID == newID {
			return true
		}
	}
	return false
}

// Check if it is the turn of the submitted name in the submitted game
func PlayersTurn(game dudo.GameState, playerName string) bool {
	for i, player := range game.Players {
		if player.Name == playerName && i == game.CurrentPlayer { // TODO change to use player ID not name and check if i = 0 initally.
			return true
		}
	}
	return false
}
