package teams

import "fmt"

type ErrEmailOrPasswdWrong struct {
	Email string
}

type ErrAlreadyLogined struct {
}

func (e *ErrEmailOrPasswdWrong) Error() string {
	return fmt.Sprintf("The email or password is wrong for Microsoft account: email: %s\n", e.Email)
}

func (e *ErrAlreadyLogined) Error() string {
	return fmt.Sprint("The microsoft accout is already logged in")
}

func NewErrEmailOrPasswdWrong(email string) *ErrEmailOrPasswdWrong {
	return &ErrEmailOrPasswdWrong{
		Email: email,
	}
}

func NewErrAlreadyLogined() *ErrAlreadyLogined {
	return &ErrAlreadyLogined{}
}
