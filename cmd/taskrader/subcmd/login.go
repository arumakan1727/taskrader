package subcmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/arumakan1727/taskrader/pkg/clients/edstem"
	"github.com/arumakan1727/taskrader/pkg/clients/gakujo"
	"github.com/arumakan1727/taskrader/pkg/clients/teams"
	"github.com/arumakan1727/taskrader/pkg/config"
	"github.com/arumakan1727/taskrader/pkg/cred"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func newLoginCmd() *cobra.Command {
	flags := struct {
		status bool
	}{}

	loginCmd := &cobra.Command{
		Use:   "login",
		Short: "ログインのための認証情報を taskrader に登録します",

		Long: "ログインのための認証情報を対話形式で taskrader に登録します。\n" +
			"\n" +
			"課題情報を取得するには、最初にこのコマンドを使って認証情報を登録しておく必要があります。\n" +
			"認証情報は手元のPCに暗号化されて保存されます。",

		ValidArgs: []string{"gakujo", "edstem", "teams"},
		Example: `
  $ taskrader login gakujo
	-> 学情へログインするためのユーザ名とパスワードを対話形式で登録します。
	   登録直後にログインを試行し、成功した場合は認証情報を暗号化して保存します。

  $ taskrader login --status
	-> 学情, EdStem, Teams それぞれについて認証情報が登録保存されているか表示します。
`,
		Args: func(cmd *cobra.Command, args []string) error {
			// --status が指定された場合はエラーチェックしない
			if flags.status {
				return nil
			}
			if len(args) < 1 {
				return fmt.Errorf("ログイン先 (%s) を指定してください", strings.Join(cmd.ValidArgs, "|"))
			}
			switch strings.ToLower(args[0]) {
			case "gakujo":
			case "edstem":
			case "teams":
				break
			default:
				return fmt.Errorf("'%s' は無効なログイン先です。指定可能なパラメータ: (%s)",
					args[0], strings.Join(cmd.ValidArgs, "|"))
			}
			if len(args) > 1 {
				return fmt.Errorf("ログイン先は1度に1つだけ指定してください")
			}
			return nil
		},

		Run: func(cmd *cobra.Command, args []string) {
			if flags.status {
				printLoginStatus()
				return
			}
			target := strings.ToLower(args[0])
			if err := interactiveLogin(target); err != nil {
				fmt.Fprintf(os.Stderr, color.RedString("Error: %s\n", err))
				os.Exit(1)
			}
			printLoginStatus()
		},
	}

	loginCmd.Flags().BoolVarP(&flags.status, "status", "s", false, "認証情報の登録状況を表示します")

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

}

func interactiveLogin(target string) error {
	t, restoreTerminalState, err := newTerminal()
	if err != nil {
		return err
	}
	defer restoreTerminalState()

	credPath, err := config.TaskraderCredentialPath()
	if err != nil {
		return err
	}
	auth := cred.LoadFromFileOrEmpty(credPath)

	switch target {
	case "gakujo":
		gakujoCred, err := interactiveLoginGakujo(t)
		if err != nil {
			return err
		}
		auth.Gakujo = *gakujoCred

	case "edstem":
		edstemCred, err := interactiveLoginEdStem(t)
		if err != nil {
			return err
		}
		auth.EdStem = *edstemCred

	case "teams":
		teamsCred, err := interactiveLoginTeams(t)
		if err != nil {
			return err
		}
		auth.Teams = *teamsCred
	}

	fmt.Fprintln(t, "認証情報を手元の環境に暗号化して保存します...")
	if err := auth.SaveToFile(credPath); err != nil {
		return fmt.Errorf("認証情報の暗号化保存に失敗しました: %v", err)
	}

	fmt.Fprintln(t, string(t.Escape.Green)+"[OK] 認証情報を暗号化して保存しました"+string(t.Escape.Reset))

	return nil
}

func interactiveLoginGakujo(t *term.Terminal) (*cred.Gakujo, error) {
	fmt.Fprint(t, "学情へログインするためのユーザ名とパスワードを入力してください。\n入力中のパスワードは表示されません。\n\n")

	username, err := askTextWithColor(t, "username: ")
	if err != nil {
		return nil, err
	}
	password, err := askPasswordWithColor(t, "password: ")
	if err != nil {
		return nil, err
	}

	fmt.Fprint(t, "入力された認証情報で学情へのログインを試みています...\n")

	err = gakujo.NewClient().Login(username, password)
	if err != nil {
		switch err := err.(type) {
		case *gakujo.ErrUsernameOrPasswdWrong:
			return nil, fmt.Errorf("ログイン失敗: おそらくユーザ名またはパスワードを間違えています (username: '%s')", err.Username)
		default:
			return nil, err
		}
	}

	fmt.Fprint(t, "ログインに成功しました。\n")

	return &cred.Gakujo{
		Username: username,
		Password: password,
	}, nil
}

func interactiveLoginEdStem(t *term.Terminal) (*cred.EdStem, error) {
	fmt.Fprint(t, "EdStem へログインするためのメールアドレスとパスワードを入力してください。\n入力中のパスワードは表示されません。\n\n")

	email, err := askTextWithColor(t, "email: ")
	if err != nil {
		return nil, err
	}
	password, err := askPasswordWithColor(t, "password: ")
	if err != nil {
		return nil, err
	}

	fmt.Fprint(t, "入力された認証情報で EdStem へのログインを試みています...\n")

	err = edstem.NewClient().Login(email, password)
	if err != nil {
		switch err := err.(type) {
		case *edstem.ErrEmailOrPasswdWrong:
			return nil, fmt.Errorf("ログイン失敗: おそらくメールアドレスまたはパスワードを間違えています (email: '%s')", err.Email)
		default:
			return nil, err
		}
	}

	fmt.Fprint(t, "ログインに成功しました。\n")

	return &cred.EdStem{
		Email:    email,
		Password: password,
	}, nil
}

func interactiveLoginTeams(t *term.Terminal) (*cred.Teams, error) {
	fmt.Fprint(t, "Teams へログインするためのメールアドレスとパスワードを入力してください。\n入力中のパスワードは表示されません。\n\n")

	email, err := askTextWithColor(t, "email: ")
	if err != nil {
		return nil, err
	}
	password, err := askPasswordWithColor(t, "password: ")
	if err != nil {
		return nil, err
	}

	fmt.Fprint(t, "入力された認証情報で Teams へのログインを試みています...\n")

	teams.ClearCookies()
	err = teams.Login(email, password, log.New(io.Discard, "", 0))
	if err != nil {
		switch err := err.(type) {
		case *teams.ErrEmailOrPasswdWrong:
			return nil, fmt.Errorf("ログイン失敗: おそらくメールアドレスまたはパスワードを間違えています (email: '%s')", err.Email)
		default:
			return nil, err
		}
	}

	fmt.Fprint(t, "ログインに成功しました。\n")

	return &cred.Teams{
		Email:    email,
		Password: password,
	}, nil
}

type funcTerminalStateRestorer = func()

// 現在の stdin, stdout を低レベルなターミナル化する (ctrl+a や ctrl+d 等のターミナルエスケープが使用可能になる)
func newTerminal() (*term.Terminal, funcTerminalStateRestorer, error) {
	stdinFd := int(os.Stdin.Fd())
	stdoutFd := int(os.Stdout.Fd())
	if !term.IsTerminal(stdinFd) || !term.IsTerminal(stdoutFd) {
		return nil, nil, fmt.Errorf("stdin,stdout should be terminal; redirects or pipes are inappropriate")
	}

	oldState, err := term.MakeRaw(stdinFd)
	if err != nil {
		return nil, nil, err
	}
	restoreFunc := func() {
		term.Restore(stdinFd, oldState)
	}

	screen := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}

	terminal := term.NewTerminal(screen, "")

	return terminal, restoreFunc, nil
}

func askTextWithColor(t *term.Terminal, prompt string) (string, error) {
	t.SetPrompt(string(t.Escape.Yellow) + prompt + string(t.Escape.Reset))
	bytes, err := t.ReadLine()
	if err != nil {
		return "", err
	}
	return string(bytes), err
}

func askPasswordWithColor(t *term.Terminal, prompt string) (string, error) {
	bytes, err := t.ReadPassword(string(t.Escape.Yellow) + prompt + string(t.Escape.Reset))
	if err != nil {
		return "", err
	}
	return string(bytes), err
}
