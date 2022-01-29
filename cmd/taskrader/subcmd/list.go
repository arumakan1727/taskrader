package subcmd

import (
	"fmt"
	"os"

	"github.com/arumakan1727/taskrader/pkg/assignment"
	"github.com/arumakan1727/taskrader/pkg/config"
	"github.com/arumakan1727/taskrader/pkg/cred"
	"github.com/arumakan1727/taskrader/pkg/view"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "統合した課題の一覧をリスト表示します。",
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if err := runListCmd(cmd, args); err != nil {
				color.Red("エラー: %s\n", err)
				os.Exit(1)
			}
		},
	}
	return listCmd
}

func runListCmd(cmd *cobra.Command, args []string) error {
	credPath, err := config.TaskraderCredentialPath()
	if err != nil {
		return err
	}
	auth := cred.LoadFromFileOrEmpty(credPath)

	color.Blue("課題を取得中...\n")
	ass, errs := assignment.FetchAll(auth)

	fmt.Printf("\n %d件の未提出課題: (締切の近い順)\n", len(ass))
	view.SortAssignments(ass)
	view.Show(ass, os.Stdout)

	if len(errs) == 0 {
		return nil
	}

	fmt.Printf("\n%d件のエラー:\n", len(errs))
	for _, err := range errs {
		color.Yellow("[%s] %s\n", err.Origin, err.Err)
	}
	return nil
}
