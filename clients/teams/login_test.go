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
	err := teams.Login(credential.Teams.Email, credential.Teams.Password, log.Default())
	if err != nil {
		t.Fatal(err)
	}
}
