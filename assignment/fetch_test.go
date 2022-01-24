package assignment_test

import (
	"testing"

	"github.com/arumakan1727/taskrader/assignment"
	"github.com/arumakan1727/taskrader/cred"
	"github.com/joho/godotenv"
)

func TestFetchAll(t *testing.T) {
	_ = godotenv.Load("../.env")

	cred := cred.LoadFromEnv()
	cred.AbortIfEmptyFieldExists()

	ass, errs := assignment.FetchAll(cred)

	for i, a := range ass {
		t.Logf("#%02d %s\n", i+1, a)
	}

	t.Logf("エラー: %d件\n", len(errs))
	for _, e := range errs {
		t.Logf("%s の課題取得でエラー: %s\n", e.Origin, e.Err)
	}
}
