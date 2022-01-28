package view_test

import (
	"os"
	"testing"
	"time"

	"github.com/arumakan1727/taskrader/pkg/assignment"
	"github.com/arumakan1727/taskrader/pkg/view"
)

func dummyAssignments() []*assignment.Assignment {
	return []*assignment.Assignment{
		{
			Origin:   assignment.OrigGakujo,
			Title:    "第11回 (12/22) お題",
			Course:   "コンパイラ",
			Deadline: time.Date(2022, time.January, 11, 24, 00, 0, 0, time.Local),
		},
		{
			Origin:   assignment.OrigEdStem,
			Title:    "当日課題(小レポート7)",
			Course:   "データベースシステム論",
			Deadline: assignment.UnknownDeadline(),
		},
		{
			Origin:   assignment.OrigTeams,
			Title:    "文字幅てすと① ↑△ ",
			Course:   "科目1",
			Deadline: time.Date(2022, time.February, 3, 23, 59, 0, 0, time.Local),
		},
		{
			Origin:   assignment.OrigTeams,
			Title:    "Ⅱ  (´・ω・｀)⚡ÄÖж",
			Course:   "科目998233353",
			Deadline: time.Date(2024, time.January, 30, 9, 59, 59, 0, time.Local),
		},
	}
}

func TestShow(t *testing.T) {
	ass := dummyAssignments()
	view.Show(ass, os.Stdout)
}
