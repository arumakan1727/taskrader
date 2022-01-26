package gakujo_test

import (
	"testing"

	"github.com/arumakan1727/taskrader/pkg/clients/gakujo"
	"github.com/arumakan1727/taskrader/pkg/cred"
	"github.com/joho/godotenv"
)

var gakujoCred cred.Gakujo

func TestMain(m *testing.M) {
	_ = godotenv.Load("../../../.env")
	c := cred.LoadFromEnv()
	c.AbortIfEmptyFieldExists()

	gakujoCred = c.Gakujo

	m.Run()
}

func TestLoginWithCorrectCred(t *testing.T) {
	c := gakujo.NewClient()
	if err := c.Login(gakujoCred.Username, gakujoCred.Password); err != nil {
		t.Fatal(err)
	}
}

func expectErrUsernameOrPasswdWrong(err error, username string, t *testing.T) {
	if err == nil {
		t.Error("Expect login error, but no error occured.")
	} else {
		switch err.(type) {
		case *gakujo.ErrUsernameOrPasswdWrong:
			e := err.(*gakujo.ErrUsernameOrPasswdWrong)
			if e.Username != username {
				t.Errorf("Expected err.Username=%s, but got %s", gakujoCred.Username, e.Username)
			}
		default:
			t.Errorf("Expected err type is *ErrUsernameOrPasswdWrong, but got %v", err)
		}
	}
}

func TestLoginWithWrongCred(t *testing.T) {
	c := gakujo.NewClient()

	t.Run("Login with correct username and wrong password", func(t *testing.T) {
		err := c.Login(gakujoCred.Username, "wrong-password")
		expectErrUsernameOrPasswdWrong(err, gakujoCred.Username, t)
	})

	t.Run("Login with wrong username and wrong password", func(t *testing.T) {
		err := c.Login("wrong-username", "wrong-password")
		expectErrUsernameOrPasswdWrong(err, "wrong-username", t)
	})
}
