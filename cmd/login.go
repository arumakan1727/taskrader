package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:     "login",
	Aliases: []string{"auth"},
	Short:   "システムへログインします。",
	Long:    "対話形式で学情やTeams等への認証情報を入力します。",

	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "gakujo":
		case "edstem":
		case "teams":
		default:
			fmt.Fprintln(os.Stderr, "不正なパラメータです: "+args[0])
		}
	},
}
