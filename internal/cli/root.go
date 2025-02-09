package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{

	Use:   "vmtoy",
	Short: "vmtoy is a CLI tool",
	Long:  `vmtoy is a CLI tool designed to help.`,
}

func init() {
	fmt.Println("init")

	rootCmd.AddCommand(StartCmd())
}

func Start() {
	rootCmd.Execute()
}
