package cli

import (
	"bias/pkg/devstreams"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var DevStreamsCommand = &cobra.Command{
	Use:   "devstreams",
	Short: "Manage the DevStreams tool.",
}

func collectStreamsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "collect",
		Short: "Start Stream Collection",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Collecting Streams")
			clientID := viper.GetString("client_id")
			clientSecret := viper.GetString("client_secret")
			devstreams.CollectStreams(clientID, clientSecret)
		},
	}
}

func configureCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "configure",
		Short: "Configure DevStreams tool",
		Run: func(cmd *cobra.Command, args []string) {
			var clientID, clientSecret string
			fmt.Print("Enter your Twitch client ID: ")
			_, err := fmt.Scan(&clientID)
			if err != nil {
				return
			}
			fmt.Print("Enter your Twitch client secret: ")
			_, err = fmt.Scan(&clientSecret)
			if err != nil {
				return
			}

			viper.Set("client_id", clientID)
			viper.Set("client_secret", clientSecret)
			err = viper.WriteConfig()
			if err != nil {
				return
			}
			fmt.Println("Configuration saved successfully!")
		},
	}

	viewCmd := &cobra.Command{
		Use:   "view",
		Short: "View the configuration file path",
		Run: func(cmd *cobra.Command, args []string) {
			configPath := viper.ConfigFileUsed()
			if configPath == "" {
				fmt.Println("No configuration file found.")
			} else {
				fmt.Println("Configuration file path:", configPath)
			}
		},
	}

	cmd.AddCommand(viewCmd)
	return cmd
}

func init() {
	cobra.OnInitialize(initConfig)

	DevStreamsCommand.AddCommand(collectStreamsCommand())
	DevStreamsCommand.AddCommand(configureCommand())
}

func initConfig() {
	viper.AddConfigPath("./cli")
	viper.SetConfigName("config")

	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("No configuration file found. Creating a new one.")
		err := viper.WriteConfig()
		if err != nil {
			fmt.Println("Error creating configuration file..")
			return
		}
	}
}
