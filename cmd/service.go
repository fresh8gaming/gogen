package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	ServiceName string
)

func GetServiceCmd() (*cobra.Command, error) {
	var cmdService = &cobra.Command{
		Use:   "service",
		Short: "",
		Long:  "",
		Run:   runServiceCmd(),
	}

	var cmdHTTP = &cobra.Command{
		Use:   "http",
		Short: "",
		Long:  "",
		Run:   runHTTPCmd(),
	}

	var cmdGRPC = &cobra.Command{
		Use:   "grpc",
		Short: "",
		Long:  "",
		Run:   runGRPCCmd(),
	}

	var cmdCron = &cobra.Command{
		Use:   "cron",
		Short: "",
		Long:  "",
		Run:   runServiceCmd(),
	}

	var err error

	cmdService.PersistentFlags().StringVarP(&Org, "org", "o", "fresh8gaming", "Github org for the monorepo (defaults to fresh8gaming)")

	cmdService.PersistentFlags().StringVarP(&ServiceName, "name", "n", "", "Name of the service")
	err = cmdService.MarkPersistentFlagRequired("name")
	if err != nil {
		return nil, err
	}

	cmdService.AddCommand(cmdHTTP, cmdGRPC, cmdCron)

	return cmdService, nil
}

func runServiceCmd() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		fmt.Print(red("Not implemented!\n"))
	}
}
