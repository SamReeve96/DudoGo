package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/SamReeve96/DudoGo/backend/gameManager"
)

func Serve() {
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

func gameHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Game created!")
	startGameFromrequest()
}

func startGameFromrequest() {
	fmt.Printf("Server: Adding game to slice \n")
	gameManager.NewGame()
}
