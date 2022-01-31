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
		Short:   "未提出課題の一覧を表示します",
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
	if errs := auth.CheckEmptyField(); len(errs) == 3 {
		fmt.Println("\n認証情報が未登録です。\n" +
			"taskrader login コマンドを使って、学情・EdStem・Teamsにログインするためのパスワード等を登録してください。\n")
	}

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
