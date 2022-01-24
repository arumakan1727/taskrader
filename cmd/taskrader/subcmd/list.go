package subcmd

import "github.com/spf13/cobra"

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "統合した課題の一覧をリスト表示します。",
	Long:    "長い説明: (いつか書く)",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
