package main

import (
	"fmt"
	"log"

	"github.com/fresh8gaming/gogen/cmd"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "nil"
	date    = "nil"
	builtBy = "dev"
)

func main() {
	cmdRepo, err := cmd.GetRepoCmd()
	if err != nil {
		log.Fatal(err)
	}

	cmdService, err := cmd.GetServiceCmd()
	if err != nil {
		log.Fatal(err)
	}

	cmdVersion := &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "",
		Long:    "",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("version: %s\n", version)
			fmt.Printf("commit: %s\n", commit)
			fmt.Printf("date: %s\n", date)
			fmt.Printf("builtBy: %s\n", builtBy)
		},
	}

	var rootCmd = &cobra.Command{Use: "gogen"}
	rootCmd.AddCommand(cmdRepo, cmdService, cmdVersion)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
