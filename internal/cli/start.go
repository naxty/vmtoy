package cli

import (
	"fmt"

	"github.com/naxty/vmtoy/internal/manager"
	"github.com/spf13/cobra"
)

var name string

func StartCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "start",
		Short: "Start a VM",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("VM name:", name)

			m := manager.NewManager()
			m.List()

			if !m.Exists(name) {
				fmt.Println("VM does not exist")
				return
			} else {
				fmt.Println("VM exists")
			}

		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Name of the VM")
	return cmd
}
