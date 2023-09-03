package cli

import (
	http "bias/http/core"
	"fmt"
	"github.com/spf13/cobra"
)

var HttpCommand = &cobra.Command{
	Use:   "http",
	Short: "Manage the Bias HTTP Server",
}

func startHttpServerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start HTTP Server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Starting HTTP Server")
			http.StartServer()
		},
	}
}

func stopHttpServerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stop HTTP Server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Stopping HTTP Server")
			http.StopServer()
		},
	}
}

func init() {
	HttpCommand.AddCommand(startHttpServerCommand())
	HttpCommand.AddCommand(stopHttpServerCommand())
}
