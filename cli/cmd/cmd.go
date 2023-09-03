package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bias",
	Short: "Bias is a CLI tool for managing HTTP server",
	Long:  `A longer description of your application goes here.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func AddCommands(commands ...*cobra.Command) error {
	for _, cmd := range commands {
		if cmd == nil {
			return fmt.Errorf("received nil command")
		}
		rootCmd.AddCommand(cmd)
	}
	return nil
}
