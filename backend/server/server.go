package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/SamReeve96/DudoGo/backend/dudo"
)

func Serve() {
	buildApp()

	cwd, err := os.Getwd()
	if err != nil {
		panic("blam!")
	}
	fmt.Printf("cwd: %s \n", cwd)

	fileServer := http.FileServer(http.Dir("../frontend/dudogo/build/"))
	http.Handle("/", fileServer)
	http.HandleFunc("/StartGame", gameHandler)

	fmt.Printf("Starting server at port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func buildApp() {
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

func gameHandler(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Path != "/hello" {
	// 	http.Error(w, "404 not found.", http.StatusNotFound)
	// 	return
	// }

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Hello!")
	startGameFromrequest()
}

func startGameFromrequest() {
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
	// this only runs the game being setup. would be better to call run game from game state (create a game interface?)
	dudo.RunGame()
}
