package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

const colorReset = "\033[0m"
const colorGreen = "\033[32m"
const colorBlue = "\033[34m"
const colorCyan = "\033[36m"

const configDirectoryName = ".spark"
const configFile = "config.yml"

// Config struct for the yaml config file
type Config struct {
	Files struct {
		IdeaDirectory string `yaml:"idea_directory"`
		IdeaFile      string `yaml:"idea_file"`
	} `yaml:"files"`
}

// getConfig returns the config from the configPath yml
func getConfig(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func main() {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	configDirectoryPath := user.HomeDir + "/" + configDirectoryName + "/"

	// creates the config file if not present
	if _, err := os.Stat(configDirectoryPath + configFile); os.IsNotExist(err) {
		initConfigFile(configDirectoryPath, configFile)
	}

	// gets the config and parses it
	config, err := getConfig(configDirectoryPath + configFile)
	if err != nil {
		log.Fatal(err)
	}

	// creates the file containing the ideas if not present
	if _, err := os.Stat(config.Files.IdeaDirectory + config.Files.IdeaFile); os.IsNotExist(err) {
		initIdeaFile(config.Files.IdeaDirectory, config.Files.IdeaFile)
	}

	// gets the new idea and pushes it
	ideaTitle, ideaDescription := getIdeaInput()
	addIdeaToFile(ideaTitle, ideaDescription, config.Files.IdeaDirectory, config.Files.IdeaFile)
	commitAndPush(ideaTitle, config.Files.IdeaDirectory)

}

// getIdeaInput returns the idea provided by the user
func getIdeaInput() (string, string) {
	ideaTitle := GetUserInput(string(colorGreen) + "Idea title: " + string(colorReset))
	ideaDescription := GetUserInput(string(colorCyan) + "Description: " + string(colorReset))

	return ideaTitle, ideaDescription
}

// initIdeaFile creates the idea file & directory if not present on the system
func initIdeaFile(ideaDirectory string, ideaFile string) {
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

// initConfigFile creates the config file with the configuration provided by the user
func initConfigFile(configDirectoryPath string, configFile string) {

	// creates the config directory if it doesn't exist
	if _, err := os.Stat(configDirectoryPath); os.IsNotExist(err) {
		os.Mkdir(configDirectoryPath, os.ModePerm)
	}

	ideaDirectory := GetUserInput(string(colorBlue)+"Enter the absolute path to your idea directory: "+string(colorReset)) + "/"
	ideaFile := "README.md"

	f, err := os.OpenFile(configDirectoryPath+configFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf(
		`files:
  idea_directory: %s
  idea_file: %s`,
		ideaDirectory, ideaFile))

	if err != nil {
		log.Fatal(err)
	}

}

// addIdeaToFile adds an idea to the idea file
func addIdeaToFile(ideaTitle string, ideaDescription string, ideaDirectory string, ideaFile string) {
	f, err := os.OpenFile(ideaDirectory+ideaFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if _, err := f.WriteString(fmt.Sprintf("## %s\n\n%s\n\n", ideaTitle, ideaDescription)); err != nil {
		log.Fatal(err)
	}
}

// commitAndPush commits the changes to the idea directory and pushes it to the git repo
func commitAndPush(message string, ideaDirectory string) {
	out, err := ExecuteCommand(ideaDirectory, "git", "add", ".")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}

	out, err = ExecuteCommand(ideaDirectory, "git", "commit", "-m", message)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}

	out, err = ExecuteCommand(ideaDirectory, "git", "push")
	if err != nil {
		// remote git not defined
		if strings.Contains(err.Error(), "exit status 128") {

			// adds the remote git
			remoteGit := GetUserInput(string(colorBlue) + "Remote repo: " + string(colorReset))
			ExecuteCommand(ideaDirectory, "git", "remote", "add", "origin", remoteGit)

			// retries to push on the repo after the remote has been added
			out, err = ExecuteCommand(ideaDirectory, "git", "push", "-u", "origin", "master")
			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Printf("%s", out)
			}
		} else {
			log.Fatal(err)
		}
	} else {
		fmt.Printf("%s", out)
	}
}

// ExecuteCommand executes the command "name" in the "directory".
// Arguments for the command can be passed by "args"
// Returns the output of the command
func ExecuteCommand(directory string, name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = directory
	return cmd.Output()
}

// GetUserInput prompts the user with a "prompt message" in the console and
// returns the user input
func GetUserInput(promptMessage string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(promptMessage)
	ideaTitle, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimRight(ideaTitle, "\n")
}
