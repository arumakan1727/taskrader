package view_test

import (
	"os"
	"testing"
	"time"

	"github.com/arumakan1727/taskrader/pkg/assignment"
	"github.com/arumakan1727/taskrader/pkg/view"
	"github.com/fatih/color"
)

func dummyAssignments() []*assignment.Assignment {
	now := time.Now()
	return []*assignment.Assignment{
		{
			Origin:   assignment.OrigTeams,
			Title:    "Ⅱ  (´・ω・｀)⚡ÄÖж",
			Course:   "科目998244353",
			Deadline: time.Date(now.Year()+1, time.January, 30, 9, 59, 59, 0, time.Local),
		},
		{
			Origin:   assignment.OrigEdStem,
			Title:    "当日課題(小レポート7)",
			Course:   "データベースシステム論",
			Deadline: assignment.UnknownDeadline(),
		},
		{
			Origin:   assignment.OrigGakujo,
			Title:    "最終レポート",
			Course:   "計算理論",
			Deadline: now.Add(23 * time.Hour),
		},
		{
			Origin:   assignment.OrigGakujo,
			Title:    "第99回 (12/22) お題",
			Course:   "コンパイラ",
			Deadline: now.Add(26 * time.Hour),
		},
		{
			Origin:   assignment.OrigTeams,
			Title:    "当日課題 ex99",
			Course:   "応用プログラミングΩ ",
			Deadline: now.Add(-1 * time.Hour),
		},
		{
			Origin:   assignment.OrigTeams,
			Title:    "文字幅てすと① ↑△ ",
			Course:   "科目1",
			Deadline: time.Date(2022, time.February, 3, 24, 01, 0, 0, time.Local),
		},
	}
}

func TestShow(t *testing.T) {
	color.NoColor = false
	ass := dummyAssignments()
	view.SortAssignments(ass)
	view.Show(ass, os.Stdout)
}
