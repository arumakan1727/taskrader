package subcmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "taskrader",
	Short: "TaskRader (仮) は様々なプラットフォームの課題情報を一括して取得・通知するアプリケーションです。",
	Long:  "長い説明: (いつか書く)",
}

// Go 言語の仕様として、init という名前の関数はパッケージが読み込まれるときに最初に自動で実行される
func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(newLoginCmd())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, color.RedString("Error: %v", err))
		os.Exit(1)
	}
}
