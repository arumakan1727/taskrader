package subcmd

import (
	"fmt"
	"os"

	"github.com/arumakan1727/taskrader/pkg/config"
	"github.com/arumakan1727/taskrader/pkg/cred"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func newStatusCmd() *cobra.Command {
	loginCmd := &cobra.Command{
		Use:   "status",
		Short: "認証情報の登録状態やその保存先ファイルへのパスを表示します",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			printLoginStatus()
		},
	}

	return loginCmd
}

func printLoginStatus() {
	credPath, err := config.TaskraderCredentialPath()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
		os.Exit(1)
	}

	auth := cred.LoadFromFileOrEmpty(credPath)

	printStatus := func(ok bool, param1, param2 string) {
		if ok {
			color.Green("[OK] %s と %s は登録済みです", param1, param2)
		} else {
			color.Red("[NG] %s と %s が未登録です", param1, param2)
		}
	}

	fmt.Println("\n[現在のログイン状態]")
	fmt.Println("-------------------------------------------------------")

	fmt.Print(" 学情:   ")
	printStatus(auth.Gakujo.Username != "" && auth.Gakujo.Password != "", "username", "password")

	fmt.Print(" EdStem: ")
	printStatus(auth.EdStem.Email != "" && auth.EdStem.Password != "", "email", "password")

	fmt.Print(" Teams:  ")
	printStatus(auth.Teams.Email != "" && auth.Teams.Password != "", "email", "password")

	fmt.Printf("\nSave path: %s\n", credPath)
}
