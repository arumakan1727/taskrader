package subcmd

import (
	"fmt"
	"os"

	"github.com/arumakan1727/taskrader/pkg/assignment"
	"github.com/arumakan1727/taskrader/pkg/server"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/toqueteos/webbrowser"
)

func newGuiCmd() *cobra.Command {
	options := struct {
		port int
	}{}

	guiCmd := &cobra.Command{
		Use:     "gui",
		Aliases: []string{"open"},
		Short:   "Webブラウザを起動してGUIで操作します",
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			gin.SetMode(gin.ReleaseMode)
			r := server.NewEngine(assignment.FetchAll)
			server.AddAssetsRoute(r)

			finished := make(chan bool, 1)

			listeningAddr := fmt.Sprintf("localhost:%d", options.port)
			go func() {
				fmt.Fprintln(os.Stderr, r.Run(listeningAddr))
				finished <- true
			}()

			if err := webbrowser.Open("http://" + listeningAddr); err != nil {
				fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
				os.Exit(1)
			}

			<-finished
		},
	}

	guiCmd.Flags().IntVarP(&options.port, "port", "p", 8777, "起動するサーバのポート番号を指定します (デフォルト: 8777)")
	return guiCmd
}
