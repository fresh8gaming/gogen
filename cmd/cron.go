package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func runCronCmd() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		fmt.Print(red("Not implemented!\n"))
	}
}
