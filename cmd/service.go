package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	ServiceName string
	Inplay      bool
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
		Run:   runCronCmd(),
	}

	var cmdCrawlerCron = &cobra.Command{
		Use:   "crawler-cron",
		Short: "",
		Long:  "",
		Run:   runCrawlerCronCmd(),
	}

	var err error

	//cmdService.PersistentFlags().StringVarP(&Org, "org", "o", "fresh8gaming", "Github org for the monorepo (defaults to fresh8gaming)")
	cmdService.PersistentFlags().StringVarP(&Team, "team", "t", "fresh8gaming", "Github org for the monorepo (defaults to fresh8gaming)")

	cmdService.PersistentFlags().StringVarP(&ServiceName, "name", "n", "", "Name of the service")
	cmdCrawlerCron.Flags().BoolVarP(
		&Inplay,
		"inplay",
		"i",
		false,
		"Whether the service is inplay or not (only used with crawler-cron, defaults to false)",
	)
	err = cmdService.MarkPersistentFlagRequired("name")
	if err != nil {
		return nil, err
	}

	cmdService.AddCommand(cmdHTTP, cmdGRPC, cmdCron, cmdCrawlerCron)

	return cmdService, nil
}

func runServiceCmd() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		fmt.Print(red("Not implemented!\n"))
	}
}
