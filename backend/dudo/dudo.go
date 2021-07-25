package dudo

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/SamReeve96/DudoGo/backend/cli"
)

// GameState - state of the current game
type GameState struct {
	Players    []Player
	Round      int
	CurrentBet Bet
	Started    bool
	// wildOnes      bool Do ones Count as any other Value
	CurrentPlayer int
}

// Player - a player of dudo
type Player struct {
	// TODO give players and ID, then when making a request for game state,
	// if the player ID is correct, return the player info, otherwise
	// you can lookup otherplayers states and cheat
	Name               string
	RemainingDiceCount int
	Dice               []int
}

// Bet - Players bet
type Bet struct {
	Count int
	Value int
}

type NewGameDetails struct {
	Players       int
	DicePerPlayer int
	FriendlyID    string
	CreatorName   string
}

var blankBet Bet = Bet{0, 0}

// Setup the game
func SetupGame(details NewGameDetails) GameState {
	fmt.Printf("Hello! Welcome to dudo go! Before we can play we need to set a few rules \n")
	gameState := GameState{}
	gameState.Players = []Player{}
	gameState = setupPlayers(gameState, details)
	fmt.Printf("Awesome! we have %v Players! and they have %v dice each \n", len(gameState.Players), gameState.Players[0].RemainingDiceCount)
	for playerNo := 0; playerNo < len(gameState.Players); playerNo++ {
		fmt.Printf("Good luck!: %v \n", gameState.Players[playerNo].Name)
	}

	gameState.CurrentPlayer = 0
	gameState.Round = 0
	gameState.Started = false
	gameState.CurrentBet = blankBet

	return gameState
}

// Setup however many Players are participating, their Names and how many Dice everyone should have
func setupPlayers(gameState GameState, details NewGameDetails) GameState {
	//TODO - replace the cli ruleset with a rules .json!
	// Get the number of Players
	for i := 0; i < details.Players; i++ {
		var newPlayer Player
		newPlayer.Name = ""
		newPlayer.RemainingDiceCount = details.DicePerPlayer
		// Set the creators name
		if i == 0 {
			newPlayer.Name = details.CreatorName
		}
		gameState.Players = append(gameState.Players, newPlayer)
	}

	return gameState
}

// Calc the total Dice in the game currently
func GetTotalDiceCount(gameState GameState) int {
	totalDice := 0
	for i := 0; i < len(gameState.Players); i++ {
		totalDice += gameState.Players[i].RemainingDiceCount
	}
	return totalDice
}

// run the game instance until it ends (total Dice == 0)
// func RunGame() {
// 	totalDice := GetTotalDiceCount()
// 	for totalDice > 0 {
// 		gameState.Round++
// 		fmt.Printf("Round: %v \n", gameState.Round)
// 		executeRound()
// 		totalDice = GetTotalDiceCount()
// 	}
// 	fmt.Printf("Game Over! Congrats: %v", gameState.Players[0].Name)
// }

// Generate Dice Values for each of the player's Dice
func RollDice(PlayerDiceNumber int) []int {
	// add a seed so numbers are randomised and non-repeatable on playthroughs
	// needs improvement as players are getting identical dice when straight after eachother
	rand.Seed(time.Now().UnixNano())
	var Dice []int
	for i := 0; i < PlayerDiceNumber; i++ {
		// get an int between ((6 - 1) + 1) +1
		// upper - lower + 1 + lower
		Dice = append(Dice, (rand.Intn(6) + 1))
	}
	return Dice
}

// Create a bet (a suggested Value and Count of that Value)
func CreateBet(gameState GameState, betValue int, betCount int) (GameState, error) {

	newBet := Bet{Value: betValue, Count: betCount}
	fmt.Printf("Bet %v, \n", newBet)
	gameState.CurrentBet = newBet

	// check if current bet is the max possible bet, i.e. Count of Value = total die if so move to eval!
	if gameState.CurrentBet.Count == GetTotalDiceCount(gameState) {
		return gameState, errors.New("Max bet has been reached, you much accuse!")
	}

	return gameState, nil
}

// Check that the bet Value and Count are acceptable ints and are possible in the current game state
// TODO - Return what validation failed to inform the player
func VaildBet(gameState GameState, Value int, Count int) bool {
	totalDiceCount := GetTotalDiceCount(gameState)
	// is the Count of Dice valid?
	if Count > totalDiceCount || Count < 0 {
		fmt.Printf("Invalid Count")
		return false
	}

	// is the Value Possible?
	if Value > 6 || Value < 1 {
		fmt.Printf("Invalid Value")
		return false
	}

	// (if the first Round then skip) if not then evaluate bet
	if gameState.CurrentBet.Count == 0 && gameState.CurrentBet.Value == 0 {
		return true
	}

	err := checkIfBetterBet(gameState, Value, Count)
	if err == nil {
		return true
	}

	fmt.Printf("error checking if bet was better: %s", err)

	return false
}

// a better bet is one with a higher Value or Count
// Unless at 6, then the Value can be 1-5 and the Count is doubled
func checkIfBetterBet(gameState GameState, Value int, Count int) error {
	// if the current bet Value is 6 (max Value) and the new bet is not, check the Count has doubled (wild ones)
	if gameState.CurrentBet.Value == 6 && Value != 6 {
		if Count >= gameState.CurrentBet.Count*2 {
			return errors.New("bet is not better than current")
		}
	}

	higherValue := Value > gameState.CurrentBet.Value
	higherCount := Count > gameState.CurrentBet.Count

	// neither Value or Count has increased
	if !higherValue && !higherCount {
		return errors.New("bet is not better than current")
	}

	return nil
}

// Check if the current bet is valid or not, then remove a die from whoever was wrong
func Accuse(gameState GameState) (GameState, error) {
	// first player of a round cant accuse, reject move (TODO: need to improve when iterative starter)
	if gameState.CurrentBet.Count == blankBet.Count &&
		gameState.CurrentBet.Value == blankBet.Value {
		return gameState, errors.New("needs to be a bet to accuse")
	}

	accusingPlayer := gameState.CurrentPlayer
	// check the last Players bet and see if it was correct or not
	accusedPlayer := getPreviousPlayer(gameState, accusingPlayer)

	// Count how many instances of the Value there currently is
	actualValueCount := getValueCount(gameState, gameState.CurrentBet.Value)

	// TODO - Print the accused player's Name
	fmt.Printf("Player bets that there were %v %v's and there are actually: %v of them\n", gameState.CurrentBet.Count, gameState.CurrentBet.Value, actualValueCount)

	// if the actual Count of the Value is less than the bet Count
	if actualValueCount >= gameState.CurrentBet.Count {
		fmt.Printf("Accuser was wrong, they lose a die \n")
		gameState.Players[accusingPlayer].RemainingDiceCount--

		if gameState.Players[accusingPlayer].RemainingDiceCount < 1 {
			removePlayer(gameState, accusingPlayer)
		}

	} else {
		fmt.Printf("Accused was wrong, they lose a die \n")
		gameState.Players[accusedPlayer].RemainingDiceCount--

		if gameState.Players[accusedPlayer].RemainingDiceCount < 1 {
			removePlayer(gameState, accusedPlayer)
		}
	}

	return gameState, nil
}

// Get the player that just made a bet
func getPreviousPlayer(gameState GameState, CurrentPlayerNo int) int {
	// if it is the first player in the slice
	if CurrentPlayerNo == 0 {
		return len(gameState.Players) - 1
	} else {
		return CurrentPlayerNo - 1
	}
}

func GetNextPlayer(gameState GameState) int {
	if (gameState.CurrentPlayer + 1) == len(gameState.Players) {
		return 0
	} else {
		return gameState.CurrentPlayer + 1
	}
}

// remove player from Players in current game
func removePlayer(gameState GameState, playerIndex int) GameState {
	fmt.Printf("%s has no more Dice, theeeeeeeey're out!", gameState.Players[playerIndex].Name)

	copy(gameState.Players[playerIndex:], gameState.Players[playerIndex+1:]) // Shift a[i+1:] left one index.
	gameState.Players[len(gameState.Players)-1] = Player{}                   // Erase last element (write empty Value).
	gameState.Players = gameState.Players[:len(gameState.Players)-1]         // Truncate slice.

	return gameState
}

// get the total Count for a provided Dice Value return true if there are that many of that Dice
func getValueCount(gameState GameState, Value int) int {
	ValueCount := 0
	for i := 0; i < len(gameState.Players); i++ {
		player := gameState.Players[i]
		// for each player go through there Dice and add 1 for each instance of the Value
		for j := 0; j < player.RemainingDiceCount; j++ {
			if player.Dice[j] == Value {
				ValueCount++
			}
		}
	}
	return ValueCount
}

func NewRound(gameState GameState) GameState {
	for i := 0; i < len(gameState.Players); i++ {
		gameState.Players[i].Dice = RollDice(gameState.Players[i].RemainingDiceCount)
	}
	gameState.CurrentBet = Bet{0, 0}
	gameState.CurrentPlayer = 0
	gameState.Round++

	return gameState
}

// Run though all the actions of a single Round of dudo (Players bet, until an accusation or max Dice Count)
// TODO - The starting player should change, either be the one that lost a Dice the previous Round or one after them if eliminated
func executeRound(gameState GameState) GameState {
	var betValue int
	var betCount int
	var err error
	var validBet bool

	NewRound(gameState)

	var RoundActive bool = true
Round:
	for RoundActive {
		gameState.CurrentPlayer = GetNextPlayer(gameState)
		if len(gameState.Players) == 1 {
			gameState.Players[0].RemainingDiceCount = 0
			// Game over
			break
		} else {
			fmt.Printf("It's %v's turn! \n", gameState.Players[gameState.CurrentPlayer].Name)
			fmt.Printf("You rolled: %v \n", gameState.Players[gameState.CurrentPlayer].Dice)
			firstBetOfRound := (gameState.CurrentBet.Count == 0 && gameState.CurrentBet.Value == 0)
			if !firstBetOfRound {
				// ask to bet or call BS
				validChoice := false
				for !validChoice {
					// TODO - remove option if you cant bet (max quantity of Dice in current bet)
					playerAction := cli.HandleInput("Do you want to Bet (B) or call BS (C)? \n")
					if playerAction == "B" || playerAction == "b" {
						validChoice = true
						validBet = false

						for !validBet {
							betValue, err = strconv.Atoi(cli.HandleInput("What Value are you betting? (1-6)\n"))
							if err != nil {
								fmt.Println(err)
							}
							totalDiceCount := GetTotalDiceCount(gameState)
							betCount, err = strconv.Atoi(cli.HandleInput(fmt.Sprintf("and how many %v 's are you betting? (1-%v) \n", betValue, totalDiceCount)))
							if err != nil {
								fmt.Println(err)
							}
							validBet = VaildBet(gameState, betValue, betCount)
							if !validBet {
								fmt.Print("Bet invalid, try again \n")
							}
						}
						CreateBet(gameState, betValue, betCount)
					} else if playerAction == "C" || playerAction == "c" { //Should cast to lowwer case but cba for now
						validChoice = true
						Accuse(gameState)

						// if the bet was called end the Round, skip the rest of the Players turns, evaluate outcome and start a new Round
						// presume RoundActive is redundant here if labeled break is working...
						break Round

					} else {
						fmt.Printf("Invalid choice \n")
					}
				}
			} else {
				for !validBet {
					betValue, err = strconv.Atoi(cli.HandleInput("What Value are you betting? (1-6)\n"))
					if err != nil {
						fmt.Println(err)
					}
					totalDiceCount := GetTotalDiceCount(gameState)
					betCount, err = strconv.Atoi(cli.HandleInput(fmt.Sprintf("and how many %v 's are you betting? (1-%v) \n", betValue, totalDiceCount)))
					if err != nil {
						fmt.Println(err)
					}
					validBet = VaildBet(gameState, betValue, betCount)
					if !validBet {
						fmt.Print("Bet invalid, try again \n")
					}
				}
				CreateBet(gameState, betValue, betCount)
			}

		}
	}

	return gameState
}
