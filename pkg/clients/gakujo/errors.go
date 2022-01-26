package gakujo

import "fmt"

type ErrUsernameOrPasswdWrong struct {
	Username string
}

func (e *ErrUsernameOrPasswdWrong) Error() string {
	return fmt.Sprintf("The username or password is wrong for Gakujo: username: %s", e.Username)
}

func NewErrUsernameOrPasswdWrong(username string) *ErrUsernameOrPasswdWrong {
	return &ErrUsernameOrPasswdWrong{
		Username: username,
	}
}
