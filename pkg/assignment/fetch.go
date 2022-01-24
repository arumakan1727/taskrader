package assignment

import (
	"io"
	"log"

	"github.com/arumakan1727/taskrader/pkg/clients/edstem"
	"github.com/arumakan1727/taskrader/pkg/clients/gakujo"
	"github.com/arumakan1727/taskrader/pkg/clients/teams"
	"github.com/arumakan1727/taskrader/pkg/cred"
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

func newErr(origin Origin, err error) *Error {
	return &Error{
		Origin: origin,
		Err:    err,
	}
}

func fetchGakujo(cred *cred.Gakujo, resultChan chan []*Assignment, errChan chan *Error) {

	client := gakujo.NewClient()

	err := client.Login(cred.Username, cred.Password)

	if err != nil {

		resultChan <- nil
		errChan <- newErr(OrigGakujo, err)
		return

	}

	tasks, err := client.GetTask()
	if err != nil {

		resultChan <- nil
		errChan <- newErr(OrigGakujo, err)
		return

	}

	result := []*Assignment{}

	for _, elem := range tasks {
		task := Assignment{
			Origin:   OrigGakujo,
			Title:    elem.Title,
			Course:   elem.Course,
			Deadline: elem.Deadline,
		}
		result = append(result, &task)
	}

	resultChan <- result
	errChan <- nil
}

func fetchEdStem(cred *cred.EdStem, resultChan chan []*Assignment, errChan chan *Error) {
	client := edstem.NewClient()
	err := client.Login(cred.Email, cred.Password)
	if err != nil {
		errChan <- newErr(OrigEdStem, err)
		resultChan <- nil
		return
	}
	announcement, err := client.JsonParse()
	if err != nil {
		errChan <- newErr(OrigEdStem, err)
		resultChan <- nil
		return
	}
	result := []*Assignment{}
	for _, ano := range announcement {
		anounce := Assignment{
			Origin:   OrigEdStem,
			Title:    ano.Title,
			Course:   ano.SubjectName,
			Deadline: UnknownDeadline(),
		}
		result = append(result, &anounce)
	}
	resultChan <- result
	errChan <- nil
}

func fetchTeams(cred *cred.Teams, resultChan chan []*Assignment, errChan chan *Error) {
	ass, err := teams.FetchAssignments(log.New(io.Discard, "", 0))
	if err != nil {
		resultChan <- nil
		errChan <- newErr(OrigTeams, err)
	}

	result := make([]*Assignment, 0, len(ass))
	for _, a := range ass {
		deadline := a.Deadline
		if deadline.IsZero() {
			deadline = UnknownDeadline()
		}
		result = append(result, &Assignment{
			Origin:   OrigTeams,
			Title:    a.Title,
			Course:   a.Course,
			Deadline: deadline,
		})
	}

	resultChan <- result
	errChan <- nil
}
