package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var name string

func StartCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "start",
		Short: "Start a VM",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("VM name:", name)

		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Name of the VM")
	return cmd
}
