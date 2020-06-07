package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// ExecuteCommand executes the command "name" in the "directory".
// Arguments for the command can be passed by "args"
// Returns the output of the command
func ExecuteCommand(directory string, name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = directory
	return cmd.Output()
}

// GetUserInput prompts the user with a "prompt message" in the console and
// returns the user input.
func GetUserInput(promptMessage string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(promptMessage)
	userInput, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	return strings.TrimRight(userInput, "\n")
}

// InitIdeaFile creates the idea file & directory if not present on the system
func InitIdeaFile(ideaDirectory string, ideaFile string) {
	// creates the directory if it does not exist
	if _, err := os.Stat(ideaDirectory); os.IsNotExist(err) {
		os.Mkdir(ideaDirectory, os.ModePerm)
	}

	// git inits the repo if it's not already done
	if _, err := os.Stat(ideaDirectory + ".git"); os.IsNotExist(err) {
		ExecuteCommand(ideaDirectory, "git", "init")
	}

	// creates the idea file
	if _, err := os.Stat(ideaDirectory + ideaFile); os.IsNotExist(err) {
		f, err := os.OpenFile(ideaDirectory+ideaFile,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		// writes the header
		if _, err := f.WriteString(fmt.Sprintf("# Ideas\n\n")); err != nil {
			log.Fatal(err)
		}
	}
}
