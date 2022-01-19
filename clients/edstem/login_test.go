package edstem

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

var (
	email    string
	password string
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("please set .env file!!", err)
	}

	email = os.Getenv("email")
	password = os.Getenv("password")
}
func TestEdstam(t *testing.T) {
	c := NewClient()
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
