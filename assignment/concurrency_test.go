package assignment

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/arumakan1727/taskrader/cred"
)

// 並行処理ができているか確かめるだけ (通信は行わない)
func TestConcurrency(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	ass, errs := fetchAllConcurrency(&cred.Credential{}, dummyFetchGakujo, dummyFetchEdStem, dummyFetchTeams)

	fmt.Printf(" TestConcurrency(): errs: %v\n", errs)
	fmt.Printf(" TestConcurrency(): assignments:\n")
	for _, a := range ass {
		fmt.Printf("    %s\n", a)
	}

	if len(errs) != 1 {
		t.Errorf("length of errs was %d, expected %d\n", len(errs), 1)
	}
	if len(ass) != 4 {
		t.Errorf("length of assignments was %d, expected %d\n", len(ass), 1)
	}

	// 課題のタイトルをカウントすることで、並行処理途中に課題が重複したり欠損したりしてないことをチェック
	titleCount := map[string]int{}
	for _, a := range ass {
		titleCount[a.Title] += 1
	}
	titleOK := true
	titleOK = titleOK && titleCount["gakujoKadai1"] == 1
	titleOK = titleOK && titleCount["gakujoKadai2"] == 1
	titleOK = titleOK && titleCount["edstemKadai1"] == 1
	titleOK = titleOK && titleCount["edstemKadai2"] == 1
	if !titleOK {
		t.Errorf("the assignments contains invalid title")
	}
}

func randomSleep() {
	ms := 500 + rand.Intn(300) // 500ms ~ 800ms
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func dummyFetchGakujo(cred *cred.Gakujo, resultChan chan []*Assignment, errChan chan *Error) {
	randomSleep()
	ass := []*Assignment{
		{
			Origin:   OrigGakujo,
			Title:    "gakujoKadai1",
			Course:   "日本語表現法",
			Deadline: time.Now(),
		},
		{
			Origin:   OrigGakujo,
			Title:    "gakujoKadai2",
			Course:   "計算理論",
			Deadline: time.Now(),
		},
	}
	resultChan <- ass
	errChan <- nil
}

func dummyFetchEdStem(cred *cred.EdStem, resultChan chan []*Assignment, errChan chan *Error) {
	randomSleep()
	ass := []*Assignment{
		{
			Origin:   OrigEdStem,
			Title:    "edstemKadai1",
			Course:   "データベースシステム論",
			Deadline: time.Now(),
		},
		{
			Origin:   OrigEdStem,
			Title:    "edstemKadai2",
			Course:   "データベースシステム論",
			Deadline: time.Now(),
		},
	}
	resultChan <- ass
	errChan <- nil
}

func dummyFetchTeams(cred *cred.Teams, resultChan chan []*Assignment, errChan chan *Error) {
	randomSleep()
	resultChan <- nil
	errChan <- &Error{
		Origin: OrigTeams,
		Err:    fmt.Errorf("Sample-teams-error"),
	}
}
