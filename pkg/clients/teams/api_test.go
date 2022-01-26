package teams_test

import (
	"io"
	"log"
	"os"
	"testing"

	"github.com/arumakan1727/taskrader/pkg/clients/teams"
	"github.com/arumakan1727/taskrader/pkg/cred"
	"github.com/joho/godotenv"
)

var credential *cred.Credential

func init() {
	_ = godotenv.Load("../../../.env")
	credential = cred.LoadFromEnv()
	credential.AbortIfEmptyFieldExists()
}

func TestLogin(t *testing.T) {
	if os.Getenv("NOW_ON_CI") != "" {
		// CI の場合は本テスト関数は実行しない
		return
	}

	teams.ClearCookies()

	t.Run("Login with correct credential and cleared cookies should be success", func(t *testing.T) {
		err := teams.Login(credential.Teams.Email, credential.Teams.Password, log.Default())
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Login should return ErrAlreadyLogined", func(t *testing.T) {
		err := teams.Login(credential.Teams.Email, credential.Teams.Password, log.Default())

		switch err.(type) {
		case *teams.ErrAlreadyLogined:
			return
		default:
			t.Errorf("Expected *teams.ErrAlreadyLogined, but got: %s", err)
		}
	})

	teams.ClearCookies()

	t.Run("ClearCookies should works & Login with incorrect password should return ErrEmailOrPasswdWrong", func(t *testing.T) {
		err := teams.Login(credential.Teams.Email, "wrong-password", log.Default())

		switch err := err.(type) {
		case *teams.ErrEmailOrPasswdWrong:
			if err.Email != credential.Teams.Email {
				t.Errorf("Expected err.Email = %s, but got %s", credential.Teams.Email, err.Email)
			}
		default:
			t.Errorf("Expected *teams.ErrAlreadyLogined, but got: %s", err)
		}
	})
}

func TestFetchAssignments(t *testing.T) {
	if os.Getenv("NOW_ON_CI") != "" {
		// CI の場合はまず空ログインしておかないと何故かうまくいかない
		_ = teams.Login(credential.Teams.Email, credential.Teams.Password, log.New(io.Discard, "", 0))
	}

	err := teams.Login(credential.Teams.Email, credential.Teams.Password, log.Default())
	if err != nil {
		switch err.(type) {
		case *teams.ErrAlreadyLogined:
			break
		default:
			t.Fatal(err)
		}
	}

	ass, err := teams.FetchAssignments(log.Default())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("len(ass) = %d", len(ass))
	for i, a := range ass {
		t.Logf("[%02d] title=%q\n  course=%q, deadline=%q\n", i+1, a.Title, a.Course, a.Deadline)
		if a.Deadline.IsZero() {
			t.Errorf("Deadline is zero; Probably failed to parse dueText")
		}
	}
}
