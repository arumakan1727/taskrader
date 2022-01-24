package teams_test

import (
	"log"
	"testing"

	"github.com/arumakan1727/taskrader/clients/teams"
	"github.com/arumakan1727/taskrader/cred"
	"github.com/joho/godotenv"
)

var credential *cred.Credential

func init() {
	_ = godotenv.Load("../../.env")
	credential = cred.LoadFromEnv()
	credential.AbortIfEmptyFieldExists()
}

func TestLogin(t *testing.T) {
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

	t.Run("ClearCookies should works & Login with incorrect credential should return ErrEmailOrPasswdWrong", func(t *testing.T) {
		err := teams.Login(credential.Teams.Email, "wrong-password", log.Default())

		switch err.(type) {
		case *teams.ErrEmailOrPasswdWrong:
			return
		default:
			t.Errorf("Expected *teams.ErrAlreadyLogined, but got: %s", err)
		}
	})
}

func TestFetchAssignments(t *testing.T) {
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
