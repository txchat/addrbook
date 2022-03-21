/*Package commands implement dapp client commands*/
package commands

import (
	"github.com/spf13/cobra"
)

/*
 * 实现合约对应客户端
 */

// Cmd chat client command
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chat",
		Short: "chat command",
		Args:  cobra.MinimumNArgs(1),
	}
	cmd.AddCommand(
	//add sub command
	)
	return cmd
}
