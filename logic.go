package gameLogic

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// GameState - state of the current game
type GameState struct {
	players    []Player
	round      int
	currentBet Bet
	// wildOnes      bool // Do ones count as other values (for now, o)
	currentPlayer int
}

// Player - a player of dudo
type Player struct {
	name               string
	remainingDiceCount int
	dice               []int
}

// Bet - players bet
type Bet struct {
	count int
	value int
}

// Global game state
var gameState GameState

// Setup the game and then execute the game loo
func RunLogic() {
	fmt.Printf("Hello! Welcome to dudo go! Before we can play we need to set a few rules \n")
	setupPlayers()
	fmt.Printf("Awesome! we have %v players!: \n", strconv.Itoa(len(gameState.players)))
	for playerNo := 0; playerNo < len(gameState.players); playerNo++ {
		fmt.Printf("Good luck!: %v \n", gameState.players[playerNo].name)
	}
	//intalise the first player
	gameState.currentPlayer = 0
	gameState.round = 0
	runGame()
}

// Just a wrapper to handle requesting user input from the CLI
func handleInput(output string) string {
	fmt.Printf(output)
	var message string
	_, err := fmt.Scanln(&message)
	if err != nil {
		fmt.Println(err)
	}
	// For debugging input issues
	// fmt.Printf("user supplied: %v \n", message)
	return message
}

// Setup however many players are participating, their names and how many dice everyone should have
func setupPlayers() {
	//TODO - replace the cli ruleset with a rules .json!
	// Get the number of players
	for gameState.players == nil {
		playersInt, err := strconv.Atoi(handleInput("How many players?: \n"))
		if err != nil {
			// throw err
			fmt.Println(err)
		}
		for i := 0; i < playersInt; i++ {
			var newPlayer Player
			gameState.players = append(gameState.players, newPlayer)
		}
	}
	for playerNo := 0; playerNo < len(gameState.players); playerNo++ {
		gameState.players[playerNo].name = handleInput(fmt.Sprintf("Enter player %v's name: \n", (playerNo + 1)))
	}
	PlayerDiceCount, err := strconv.Atoi(handleInput("How many Dice should each person have?: "))
	if err != nil {
		// throw err
		fmt.Println(err)
	}
	for playerNo := 0; playerNo < len(gameState.players); playerNo++ {
		gameState.players[playerNo].remainingDiceCount = PlayerDiceCount
	}
}

// Calc the total dice in the game currently
func getTotalDiceCount() int {
	totalDice := 0
	for i := 0; i < len(gameState.players); i++ {
		totalDice += gameState.players[i].remainingDiceCount
	}
	return totalDice
}

// run the game instance until it ends (total dice == 0)
func runGame() {
	totalDice := getTotalDiceCount()
	for totalDice > 0 {
		gameState.round++
		fmt.Printf("round: %v \n", gameState.round)
		executeRound()
		totalDice = getTotalDiceCount()
	}
	fmt.Printf("Game Over! Congrats: %v", gameState.players[0].name)
}

// Generate Dice values for each of the player's dice
func rollDice(PlayerDiceNumber int) []int {
	// add a seed so numbers are randomised and non-repeatable on playthroughs
	rand.Seed(time.Now().UnixNano())
	var dice []int
	for i := 0; i < PlayerDiceNumber; i++ {
		// get an int between ((6 - 1) + 1) +1
		// upper - lower + 1 + lower
		dice = append(dice, (rand.Intn(6) + 1))
	}
	return dice
}

// Create a bet (a suggested Value and Count of that value)
func createBet() {
	validBet := false
	var betValue int
	var betCount int
	var err error
	for validBet == false {
		betValue, err = strconv.Atoi(handleInput("What value are you betting? (1-6)\n"))
		if err != nil {
			fmt.Println(err)
		}
		totalDiceCount := getTotalDiceCount()
		betCount, err = strconv.Atoi(handleInput(fmt.Sprintf("and how many %v 's are you betting? (1-%v) \n", betValue, totalDiceCount)))
		if err != nil {
			fmt.Println(err)
		}
		validBet = validateBet(betValue, betCount)
		if validBet == false {
			fmt.Print("Bet invalid, try again \n")
		}
	}

	newBet := Bet{value: betValue, count: betCount}
	fmt.Printf("Bet %v, \n", newBet)
	gameState.currentBet = newBet

	// check if current bet is the max possible bet, i.e. Count of value = total die if so move to eval!
	if gameState.currentBet.count == getTotalDiceCount() {
		evaluateCurrentBet()
	}
}

// Check that the bet value and count are acceptable ints and are possible in the current game state
// TODO - Return what validation failed to inform the player
func validateBet(value int, count int) bool {
	totalDiceCount := getTotalDiceCount()
	// is the count of Dice valid?
	if count > totalDiceCount || count < 0 {
		return false
	}

	// is the Value Possible?
	if value > 6 || value < 1 {
		return false
	}

	// (if the first round then skip) if not then evaluate bet
	if (gameState.currentBet.count == 0 && gameState.currentBet.value == 0) || checkIfBetterBet(value, count) {
		return true
	}

	return false
}

// a better bet is one with a higher value or count
// Unless at 6, then the value can be 1-5 and the count is doubled
func checkIfBetterBet(value int, count int) bool {
	// if the current bet value is 6 (max value) and the new bet is not, check the count has doubled (wild ones)
	if gameState.currentBet.value == 6 && value != 6 {
		if count >= gameState.currentBet.count*2 {
			return true
		}
	}

	higherValue := value > gameState.currentBet.value
	higherCount := count > gameState.currentBet.count

	// neither value or count has increased
	if !higherValue && !higherCount {
		return false
	}

	return true
}

// Check if the current bet is valid or not, then remove a die from whoever was wrong
func accuse() {
	accusingPlayer := gameState.currentPlayer
	// check the last players bet and see if it was correct or not
	accusedPlayer := getPreviousPlayer(accusingPlayer)

	// count how many instances of the value there currently is
	actualValueCount := getValueCount(gameState.currentBet.value)

	// TODO - Print the accused player's name
	fmt.Printf("Player bets that there were %v %v's and there are actually: %v of them\n", gameState.currentBet.count, gameState.currentBet.value, actualValueCount)

	// if the actual count of the value is less than the bet count
	if actualValueCount >= gameState.currentBet.count {
		fmt.Printf("Accuser was wrong, they lose a die \n")
		gameState.players[accusingPlayer].remainingDiceCount--

		if gameState.players[accusingPlayer].remainingDiceCount < 1 {
			removePlayer(accusingPlayer)
		}

	} else {
		fmt.Printf("Accused was wrong, they lose a die \n")
		gameState.players[accusedPlayer].remainingDiceCount--

		if gameState.players[accusedPlayer].remainingDiceCount < 1 {
			removePlayer(accusedPlayer)
		}
	}
}

// Get the player that just made a bet
func getPreviousPlayer(currentPlayerNo int) int {
	// if it is the first player in the slice
	if currentPlayerNo == 0 {
		return len(gameState.players) - 1
	} else {
		return currentPlayerNo - 1
	}
}

func getNextPlayer() {
	if (gameState.currentPlayer + 1) == len(gameState.players) {
		gameState.currentPlayer = 0
	} else {
		gameState.currentPlayer++
	}
}

// Check if Bet is possible (Occurs when the max count is reached (bet.count == total no. dice in the game))
func evaluateCurrentBet() {
	actualValueCount := getValueCount(gameState.currentBet.value)

	if actualValueCount < gameState.currentBet.count {
		fmt.Printf("last player was wrong, they lose a die \n")
		previousPlayer := getPreviousPlayer(gameState.currentPlayer)
		gameState.players[previousPlayer].remainingDiceCount--

		if gameState.players[previousPlayer].remainingDiceCount < 1 {
			removePlayer(previousPlayer)
		}

	}

}

// remove player from players in current game
func removePlayer(playerIndex int) {
	fmt.Printf("%s has no more dice, theeeeeeeey're out!", gameState.players[playerIndex].name)

	copy(gameState.players[playerIndex:], gameState.players[playerIndex+1:]) // Shift a[i+1:] left one index.
	gameState.players[len(gameState.players)-1] = Player{}                   // Erase last element (write empty value).
	gameState.players = gameState.players[:len(gameState.players)-1]         // Truncate slice.
}

// get the total count for a provided Dice value return true if there are that many of that dice
func getValueCount(value int) int {
	valueCount := 0
	for i := 0; i < len(gameState.players); i++ {
		player := gameState.players[i]
		// for each player go through there dice and add 1 for each instance of the value
		for j := 0; j < player.remainingDiceCount; j++ {
			if player.dice[j] == value {
				valueCount++
			}
		}
	}
	return valueCount
}

// Run though all the actions of a single round of dudo (Players bet, until an accusation or max dice count)
// TODO - The starting player should change, either be the one that lost a dice the previous round or one after them if eliminated
func executeRound() {

	for i := 0; i < len(gameState.players); i++ {
		gameState.players[i].dice = rollDice(gameState.players[i].remainingDiceCount)
	}
	gameState.currentBet = Bet{0, 0}
	roundActive := true
round:
	for roundActive {
		getNextPlayer()
		if len(gameState.players) == 1 {
			gameState.players[0].remainingDiceCount = 0
			// Game over
			break
		} else {
			fmt.Printf("It's %v's turn! \n", gameState.players[gameState.currentPlayer].name)
			fmt.Printf("You rolled: %v \n", gameState.players[gameState.currentPlayer].dice)
			firstBetOfRound := (gameState.currentBet.count == 0 && gameState.currentBet.value == 0)
			if !firstBetOfRound {
				// ask to bet or call BS
				validChoice := false
				for validChoice == false {
					// TODO - remove option if you cant bet (max quantity of dice in current bet)
					playerAction := handleInput("Do you want to Bet (B) or call BS (C)? \n")
					if playerAction == "B" || playerAction == "b" {
						validChoice = true
						createBet()
					} else if playerAction == "C" || playerAction == "c" { //Should cast to lowwer case but cba for now
						validChoice = true
						accuse()

						// if the bet was called end the round, skip the rest of the players turns, evaluate outcome and start a new round
						// presume roundActive is redundant here if labeled break is working...
						break round

					} else {
						fmt.Printf("Invalid choice \n")
					}
				}
			} else {
				createBet()
			}

		}
	}
}
