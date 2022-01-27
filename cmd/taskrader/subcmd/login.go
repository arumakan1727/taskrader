package subcmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/arumakan1727/taskrader/pkg/clients/gakujo"
	"github.com/arumakan1727/taskrader/pkg/config"
	"github.com/arumakan1727/taskrader/pkg/cred"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var loginCmd = &cobra.Command{
	Use:       "login",
	Aliases:   []string{"auth"},
	Short:     "システムへログインします。",
	Long:      "対話形式で学情やTeams等への認証情報を入力します。",
	ValidArgs: []string{"gakujo", "edstem", "teams"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("ログイン先 (%s) を指定してください", strings.Join(cmd.ValidArgs, "|"))
		}
		if len(args) > 1 {
			return fmt.Errorf("ログイン先は1度に1つだけ指定してください")
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
		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		target := strings.ToLower(args[0])
		if err := interactiveLogin(target); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)

		}
	},
}

func interactiveLogin(target string) error {
	t, restoreTerminalState, err := newTerminal()
	if err != nil {
		return err
	}
	defer restoreTerminalState()

	savePath, err := config.TaskraderCredentialPath()
	if err != nil {
		return err
	}
	auth := cred.LoadFromFileOrEmpty(savePath)

	switch target {
	case "gakujo":
		gakujoCred, err := interactiveLoginGakujo(t)
		if err != nil {
			return err
		}
		auth.Gakujo = *gakujoCred

	case "edstem":
	case "teams":
	}

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
		switch err.(type) {
		case *gakujo.ErrUsernameOrPasswdWrong:
			return nil, fmt.Errorf("ログイン失敗: おそらくユーザ名またはパスワードを間違えています (username: %s)", username)
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
