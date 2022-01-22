package teams_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/arumakan1727/taskrader/clients/teams"
	"github.com/arumakan1727/taskrader/cred"
	"github.com/joho/godotenv"
)

var credential *cred.Credential

func init() {
	if err := godotenv.Load("../../.env"); err != nil {
		panic("Please put ../../.env !!")
	}

	credential = cred.LoadFromEnv()
	if errs := credential.CheckEmptyField(); len(errs) > 0 {
		for _, err := range errs {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}

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
