package cmd

import (
	"fmt"
	"log"
	"os"
	"spark/utils"
	"strings"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "spark",
	Short: "Adds your idea to your list and pushes it to your git",
	Long: `
	Spark is a tool to input and store your ideas

	Spark enables you to quickly write your ideas without leaving the terminal. 
	Each addition is automatically pushed to your git repository. 
	Your ideas are stored in a README.md file for easy reading.`,

	Run: func(cmd *cobra.Command, args []string) {
		ideaDirectory := viper.GetString("idea_directory") + "/"
		ideaFile := viper.GetString("idea_file")

		// creates the file containing the ideas if not present
		if _, err := os.Stat(ideaDirectory + ideaFile); os.IsNotExist(err) {
			utils.InitIdeaFile(ideaDirectory, ideaFile)
		}

		ideaTitle, ideaDescription := getIdeaInput()
		addIdeaToFile(ideaTitle, ideaDescription, ideaDirectory, ideaFile)
		commitAndPush(ideaTitle, ideaDirectory)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file
func initConfig() {

	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Search config in home directory with name ".spark"
	viper.AddConfigPath(home)
	viper.SetConfigName(".spark")
	viper.SetConfigType("yaml")

	viper.SetDefault("idea_file", "README.md")

	// If a config file is found, read it in.
	err = viper.ReadInConfig()
	validateIdeaDirectory()

	if err != nil {
		if err := viper.SafeWriteConfigAs(home + "/.spark.yaml"); err != nil {
			if os.IsNotExist(err) {
				err = viper.WriteConfigAs(home + "/.spark.yaml")
			}
		}
	}

}

func validateIdeaDirectory() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Set if idea_directory if not set
	if viper.Get("idea_directory") == nil {
		defaultIdeaDirectory := home + "/ideas"

		ideaDirectory := utils.GetUserInput(
			fmt.Sprintf(
				"%sEnter the absolute path to your idea directory [ %s ]: %s",
				string(utils.ColorWhite),
				defaultIdeaDirectory,
				string(utils.ColorReset),
			),
		)

		if ideaDirectory == "" {
			ideaDirectory = defaultIdeaDirectory
		}
		viper.Set("idea_directory", ideaDirectory)
	}

	viper.WriteConfig()
}

// getIdeaInput returns the idea provided by the user
func getIdeaInput() (string, string) {
	ideaTitle := utils.GetUserInput(fmt.Sprintf("%sIdea title: %s", string(utils.ColorGreen), string(utils.ColorReset)))
	ideaDescription := utils.GetUserInput(fmt.Sprintf("%sDescription: %s", string(utils.ColorCyan), string(utils.ColorReset)))

	return ideaTitle, ideaDescription
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
	out, err := utils.ExecuteCommand(ideaDirectory, "git", "add", ".")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}

	out, err = utils.ExecuteCommand(ideaDirectory, "git", "commit", "-m", message)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}

	out, err = utils.ExecuteCommand(ideaDirectory, "git", "push")
	if err != nil {
		// remote git not defined
		if strings.Contains(err.Error(), "exit status 128") {

			// adds the remote git
			remoteGit := utils.GetUserInput(fmt.Sprintf("%sRemote repo: %s", string(utils.ColorWhite), string(utils.ColorReset)))
			utils.ExecuteCommand(ideaDirectory, "git", "remote", "add", "origin", remoteGit)

			// retries to push on the repo after the remote has been added
			out, err = utils.ExecuteCommand(ideaDirectory, "git", "push", "-u", "origin", "master")
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
