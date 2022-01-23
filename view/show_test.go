package view

import (
	"fmt"
	"os"
	"testing"

	"github.com/arumakan1727/taskrader/assignment"
	"github.com/arumakan1727/taskrader/cred"
	"github.com/joho/godotenv"
)

func TestShow(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		panic("Please put ../.env !")
	}

	cred := cred.LoadFromEnv()
	if errs := cred.CheckEmptyField(); len(errs) > 0 {
		for _, e := range errs {
			fmt.Fprintln(os.Stderr, "Error: "+e.Error())
		}
		os.Exit(1)
	}

	ass, errs := assignment.FetchAll(cred)

	t.Logf("エラー: %d件\n", len(errs))
	for _, e := range errs {
		t.Logf("%s の課題取得でエラー: %s\n", e.Origin, e.Err)
	}

	Show(ass)

}