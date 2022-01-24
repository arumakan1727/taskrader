package view

import (
	"testing"

	"github.com/arumakan1727/taskrader/pkg/assignment"
	"github.com/arumakan1727/taskrader/pkg/cred"
	"github.com/joho/godotenv"
)

func TestShow(t *testing.T) {
	_ = godotenv.Load("../../.env")
	cred := cred.LoadFromEnv()
	cred.AbortIfEmptyFieldExists()

	ass, errs := assignment.FetchAll(cred)

	t.Logf("エラー: %d件\n", len(errs))
	for _, e := range errs {
		t.Logf("%s の課題取得でエラー: %s\n", e.Origin, e.Err)
	}

	Show(ass)
}
