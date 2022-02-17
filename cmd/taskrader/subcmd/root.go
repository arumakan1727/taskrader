package subcmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "taskrader",
	Short: "TaskRader は学情・EdStem・Teamsのの課題情報を一括して取得するアプリケーションです。",
	Long: `
████████╗ █████╗ ███████╗██╗  ██╗██████╗  █████╗ ██████╗ ███████╗██████╗
╚══██╔══╝██╔══██╗██╔════╝██║ ██╔╝██╔══██╗██╔══██╗██╔══██╗██╔════╝██╔══██╗
   ██║   ███████║███████╗█████╔╝ ██████╔╝███████║██║  ██║█████╗  ██████╔╝
   ██║   ██╔══██║╚════██║██╔═██╗ ██╔══██╗██╔══██║██║  ██║██╔══╝  ██╔══██╗
   ██║   ██║  ██║███████║██║  ██╗██║  ██║██║  ██║██████╔╝███████╗██║  ██║
   ╚═╝   ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝╚═════╝ ╚══════╝╚═╝  ╚═╝

TaskRader は学情・EdStem・Teamsのの課題情報を一括して取得するアプリケーションです。
コマンドラインによる操作も、WebブラウザによるGUI操作もできます。

 --------------------------------------------------------------------------------

【必要なパッケージ】
  Teams の課題取得には Selenium WebDriver を使用します。
  Ubuntu をご利用なら
    sudo apt update && sudo apt install -y chromium-chromedriver
  を実行すれば WebDriver をインストールできます。

【GUIで操作】
  taskrader gui を実行すればよいです。Webブラウザが起動します。

【CLIで操作】
  コマンドラインで操作するには、2つのサブコマンド 'login', 'list' を使用します。

 --------------------------------------------------------------------------------
`,
}

// Go 言語の仕様として、init という名前の関数はパッケージが読み込まれるときに最初に自動で実行される
func init() {
	rootCmd.AddCommand(newListCmd())
	rootCmd.AddCommand(newLoginCmd())
	rootCmd.AddCommand(newStatusCmd())
	rootCmd.AddCommand(newGuiCmd())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, color.RedString("Error: %v", err))
		os.Exit(1)
	}
}
