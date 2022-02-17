package assignment_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/arumakan1727/taskrader/pkg/assignment"
	"github.com/arumakan1727/taskrader/pkg/cred"
	"github.com/joho/godotenv"
)

func TestFetchAll(t *testing.T) {
	if os.Getenv("NOW_ON_CI") != "" {
		// CI の場合は本テスト関数は実行しない
		return
	}

	_ = godotenv.Load("../../.env")

	cred := cred.LoadFromEnv()
	cred.AbortIfEmptyFieldExists()

	ass, errs := assignment.FetchAll(cred)

	for i, a := range ass {
		fmt.Printf("#%02d %s\n", i+1, a)
	}

	fmt.Printf("エラー: %d件\n", len(errs))
	for _, e := range errs {
		fmt.Printf("%s の課題取得でエラー: %s\n", e.Origin, e.Err)
	}
}
