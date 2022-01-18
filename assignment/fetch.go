package assignment

import (
	"fmt"

	"github.com/arumakan1727/taskrader/cred"
)

type gakujoFetechGoroutine = func(*cred.Gakujo, chan []*Assignment, chan *Error)
type edstemFetechGoroutine = func(*cred.EdStem, chan []*Assignment, chan *Error)
type teamsFetechGoroutine = func(*cred.Teams, chan []*Assignment, chan *Error)

func FetchAll(cred *cred.Credential) ([]*Assignment, []*Error) {
	return fetchAllConcurrency(cred, fetchGakujo, fetchEdStem, fetchTeams)
}

// ゴルーチンを用いて3つ並行で課題取得を行う
func fetchAllConcurrency(
	cred *cred.Credential,
	gakujoFetcher gakujoFetechGoroutine,
	edstemFetcher edstemFetechGoroutine,
	teamsFetcher teamsFetechGoroutine,
) ([]*Assignment, []*Error) {

	assignmentsChan := make(chan []*Assignment, 3)
	errChan := make(chan *Error, 3)
	defer close(assignmentsChan)
	defer close(errChan)

	go gakujoFetcher(&cred.Gakujo, assignmentsChan, errChan)
	go edstemFetcher(&cred.EdStem, assignmentsChan, errChan)
	go teamsFetcher(&cred.Teams, assignmentsChan, errChan)

	assignments := make([]*Assignment, 0, 32)
	errs := make([]*Error, 0, 3)

	for i := 0; i < 3; i++ {
		ass := <-assignmentsChan
		err := <-errChan
		if err != nil {
			errs = append(errs, err)
		} else {
			assignments = append(assignments, ass...)
		}
	}

	return assignments, errs
}

func fetchGakujo(cred *cred.Gakujo, resultChan chan []*Assignment, errChan chan *Error) {
	resultChan <- nil
	errChan <- &Error{
		Origin: OrigGakujo,
		Err:    fmt.Errorf("assignment.fetchGakujo が未実装です"),
	}
}

func fetchEdStem(cred *cred.EdStem, resultChan chan []*Assignment, errChan chan *Error) {
	resultChan <- nil
	errChan <- &Error{
		Origin: OrigEdStem,
		Err:    fmt.Errorf("assignment.fetchEdStem が未実装です"),
	}
}

func fetchTeams(cred *cred.Teams, resultChan chan []*Assignment, errChan chan *Error) {
	resultChan <- nil
	errChan <- &Error{
		Origin: OrigTeams,
		Err:    fmt.Errorf("assignment.fetchTeams が未実装です"),
	}
}
