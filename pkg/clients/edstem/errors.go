package edstem

import (
	"fmt"
)

type ErrEmailOrPasswdWrong struct {
	Email string
}

func (e *ErrEmailOrPasswdWrong) Error() error {
	return fmt.Errorf(fmt.Sprintf("The email or password is wrong for Edstem account: email: %s\n", e.Email))
}

func NewErrEmailOrPasswdWrong(email string) *ErrEmailOrPasswdWrong {
	return &ErrEmailOrPasswdWrong{
		Email: email,
	}
}
