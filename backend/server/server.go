package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/SamReeve96/DudoGo/backend/dudo"
	"github.com/SamReeve96/DudoGo/backend/gameManager"
)

type gameDetails struct {
	players       int
	dicePerPlayer int
}

func Serve() {
	cwd, err := os.Getwd()
	if err != nil {
		panic("blam!")
	}
	fmt.Printf("cwd: %s \n", cwd)

	fileServer := http.FileServer(http.Dir("../frontend/dudogo/build/"))
	http.Handle("/", fileServer)
	http.HandleFunc("/NewGame", newGame)
	// http.HandleFunc("/joinGame", gameHandler)

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
	fmt.Print("blam \n")
	switch r.Method {
	case "GET":
		fmt.Printf("Cant get a new game!")
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Printf("ParseForm() err: %v", err)
			return
		}
		fmt.Printf("Post from website! r.PostFrom = %v\n", r.PostForm)
		players, err := strconv.Atoi(r.FormValue("players"))
		dicePerPlayer, err := strconv.Atoi(r.FormValue("dicePerPlayer"))

		if err != nil {
			fmt.Printf("Invalid request values, players %s, dicePerPlayer %s", r.FormValue("players"), r.FormValue("dicePerPlayer"))
			return
		}

		details := dudo.NewGameDetails{
			Players:       players,
			DicePerPlayer: dicePerPlayer,
		}

		fmt.Printf("Server: Adding game to slice \n")
		gameManager.NewGame(details)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}
