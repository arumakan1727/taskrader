package edstem_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/arumakan1727/taskrader/pkg/clients/edstem"
	"github.com/arumakan1727/taskrader/pkg/cred"
	"github.com/joho/godotenv"
)

var (
	email    string
	password string
)

func init() {
	_ = godotenv.Load("../../../.env")

	cred := cred.LoadFromEnv()
	cred.AbortIfEmptyFieldExists()
	email = cred.EdStem.Email
	password = cred.EdStem.Password
}

func TestEdstam(t *testing.T) {
	c := edstem.NewClient()
	err := c.Login(email, password)
	if err != nil {
		log.Fatal(err)
	}
	announcement, err := c.JsonParse()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(announcement)
}

func TestEdstemWrongPasswd(t *testing.T) {
	c := edstem.NewClient()
	err := c.Login(email, password+"wrong_password")

	switch err := err.(type) {
	case *edstem.ErrEmailOrPasswdWrong:
		if err.Email != email {
			t.Errorf("err.Email: Expected %s, but got %s", email, err.Email)
		}
		break

	default:
		t.Errorf("Expected *edstem.ErrEmailOrPasswdWrong, but got %v", err)
	}
}
