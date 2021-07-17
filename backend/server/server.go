package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/SamReeve96/DudoGo/backend/dudo"
	"github.com/SamReeve96/DudoGo/backend/gameManager"
	"github.com/google/uuid"
)

func Serve() {
	cwd, err := os.Getwd()
	if err != nil {
		panic("blam!")
	}
	fmt.Printf("cwd: %s \n", cwd)

	fileServer := http.FileServer(http.Dir("../frontend/dudogo/build/"))
	http.Handle("/", fileServer)
	http.HandleFunc("/NewGame", newGame)
	http.HandleFunc("/JoinGame", joinGame)
	http.HandleFunc("/StartGame", startGame)
	http.HandleFunc("/GameState", gameState)

	fmt.Printf("Starting server at port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func BuildApp() {
	if runtime.GOOS == "windows" {
		// go to site dir and build
		os.Chdir("../frontend/dudogo")
		fmt.Printf("Building web app \n")
		cmd := exec.Command("npm", "run-script", "build")

		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}

		// go back to backend
		os.Chdir("../../backend/")
	} else {
		panic("unsupported OS")
	}
}

func newGame(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Printf("Cant get a new game!")
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Printf("ParseForm() err: %v", err)
			return
		}
		players, err := strconv.Atoi(r.FormValue("players"))
		dicePerPlayer, err := strconv.Atoi(r.FormValue("dicePerPlayer"))
		friendlyID := (r.FormValue("friendlyID"))
		creatorName := (r.FormValue("creatorName"))

		if err != nil {
			fmt.Printf(`Invalid request values:  
			players: %d, 
			dicePerPlayer %d
			friendlyID %s `,
				players,
				dicePerPlayer,
				friendlyID)

			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if players == 0 || dicePerPlayer == 0 || friendlyID == "" {
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		if gameManager.FriendlyIDInUse(friendlyID) {
			w.WriteHeader(http.StatusConflict)
			// improve to make response inform user that the friendly ID is in use
			return
		}

		details := dudo.NewGameDetails{
			Players:       players,
			DicePerPlayer: dicePerPlayer,
			FriendlyID:    friendlyID,
			CreatorName:   creatorName,
		}

		fmt.Printf("Server: Adding game to slice \n")
		gameManager.NewGame(details)
		w.WriteHeader(http.StatusOK)

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func joinGame(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Printf("ParseForm() err: %v", err)
			return
		}
		friendlyID := (r.FormValue("friendlyID"))
		playerName := (r.FormValue("playerName"))

		joined := gameManager.JoinGame(friendlyID, playerName)

		if joined {
			fmt.Printf("Added new player: %s to game: %s \n", playerName, friendlyID)
			w.WriteHeader(http.StatusOK)
			return
		}

		fmt.Printf("Couldn't add new player: %s to game: %s \n", playerName, friendlyID)
		w.WriteHeader(http.StatusUnauthorized)

	case "POST":
		fmt.Printf("Cant post a new game!")
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func startGame(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Printf("Cant get a start game!")
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Printf("ParseForm() err: %v", err)
			return
		}

		friendlyID := (r.FormValue("friendlyID"))
		creatorName := (r.FormValue("playerName"))

		err := gameManager.StartGame(friendlyID, creatorName)

		if err != nil {
			fmt.Print(err.Error())
			w.WriteHeader(http.StatusUnauthorized)
		}

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func gameState(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Printf("ParseForm() err: %v", err)
			return
		}

		friendlyID := (r.FormValue("friendlyID"))
		//TODO change this to a player ID so that you cant lookup competing players states

		// creatorName := (r.FormValue("playerName"))

		// for now it's fine to get and return the whole thing, but when running, remove the player data other than the current player
		gameState := gameManager.GetGame(friendlyID)

		if gameState.Id == uuid.Nil {
			fmt.Print("Game not found")
			w.WriteHeader(http.StatusNotFound)
		}

		w.Header().Set("Content-Type", "application/json")
		jsonResp, err := json.Marshal(gameState)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
		return

	case "POST":
		fmt.Printf("Cant post a gamestate!")
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}
