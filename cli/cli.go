// Basic logic for CLI interaction, will be deprecated one alternative interfaces are implemented

package cli

import "fmt"

// Just a wrapper to handle requesting user input from the CLI
func HandleInput(output string) string {
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
