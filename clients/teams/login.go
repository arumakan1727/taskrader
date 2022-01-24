package teams

import (
	"log"
	"strings"
	"time"

	"github.com/sclevine/agouti"
)

const (
	elemNameEmail       = "loginfmt"
	elemNamePassword    = "passwd"
	elemIDConfirmButton = "idSIButton9"
	elemIDPasswordError = "passwordError"
)

func Login(email, password string, logger *log.Logger) error {
	opt, err := myChromeOptions()
	if err != nil {
		return err
	}
	driver := agouti.ChromeDriver(opt)
	if err := driver.Start(); err != nil {
		return err
	}
	defer driver.Stop()

	page, err := driver.NewPage()
	if err != nil {
		return err
	}

	logger.Printf("Opening %s ...", loginURL)
	if err := page.Navigate(loginURL); err != nil {
		return err
	}
	time.Sleep(time.Second)

	{
		url, err := page.URL()
		if err != nil {
			return err
		}
		if strings.HasPrefix(url, "https://www.office.com") {
			return NewErrAlreadyLogined()
		}
	}

	logger.Printf("Sending email %s ...", email)
	if err := sendKeys(page, elemNameEmail, email, 7*time.Second); err != nil {
		return err
	}
	if err := clickButtonByID(page, elemIDConfirmButton, 3*time.Second); err != nil {
		return err
	}
	time.Sleep(time.Second)

	logger.Printf("Sending password ...")
	if err := sendKeys(page, elemNamePassword, password, 7*time.Second); err != nil {
		return err
	}
	if err := clickButtonByID(page, elemIDConfirmButton, 3*time.Second); err != nil {
		return err
	}
	time.Sleep(time.Second)

	if existsElementHavingID(page, elemIDPasswordError, 100*time.Millisecond, 300*time.Millisecond) {
		return NewErrEmailOrPasswdWrong(email)
	}
	if err := clickButtonByID(page, elemIDConfirmButton, 7*time.Second); err != nil {
		return err
	}

	logger.Printf("Login succeeded.")
	return nil
}
