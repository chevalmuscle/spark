package cmd

import (
	"fmt"
	"os"
	"spark/utils"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configurate the tool",
	Long: `
	Asks the user to fill back the configuration options.
	The file location is $HOME/.spark.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		config()

		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err != nil {
			if err := viper.SafeWriteConfigAs(home + "/.spark.yaml"); err != nil {
				if os.IsNotExist(err) {
					err = viper.WriteConfigAs(home + "/.spark.yaml")
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func config() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Sets idea_directory
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

	// Sets idea_file
	defaultIdeaFile := "README.md"

	ideaFile := utils.GetUserInput(
		fmt.Sprintf(
			"%sEnter the wanted name of your idea file [ %s ]: %s",
			string(utils.ColorWhite),
			defaultIdeaFile,
			string(utils.ColorReset),
		),
	)

	if ideaFile == "" {
		ideaFile = defaultIdeaFile
	}
	viper.Set("idea_file", ideaFile)

	viper.WriteConfig()
}
