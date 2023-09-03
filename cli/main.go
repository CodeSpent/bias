package main

import (
	cli "bias/cli/cmd"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/signal"
)

func main() {
	logger := log.New(os.Stdout, "bias: ", log.LstdFlags|log.Lshortfile)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			logger.Println("caught signal:", sig)
			os.Exit(1)
		}
	}()

	commands := []*cobra.Command{
		cli.HttpCommand,
		cli.DevStreamsCommand,
	}

	for _, cmd := range commands {
		addCommandOrExit(cmd)
	}

	if err := cli.AddCommands(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := cli.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func addCommandOrExit(cmd *cobra.Command) {
	err := cli.AddCommands(cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
