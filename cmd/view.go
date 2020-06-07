package cmd

import (
	"bufio"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "Prints the idea list using less",
	Long: `
	Prints the idea list using less from /usr/bin/less`,
	Run: func(cmd *cobra.Command, args []string) {
		ideaFilePath := viper.GetString("idea_directory") + "/" + viper.GetString("idea_file")

		command := exec.Command("/usr/bin/less")
		file, err := os.Open(ideaFilePath)

		if err != nil {
			log.Fatal(err)
		}

		command.Stdin = bufio.NewReader(file)
		command.Stdout = os.Stdout

		err = command.Run()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)
}
