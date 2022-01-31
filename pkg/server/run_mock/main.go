package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/arumakan1727/taskrader/pkg/assignment"
	"github.com/arumakan1727/taskrader/pkg/cred"
	"github.com/arumakan1727/taskrader/pkg/server"
	"github.com/gin-gonic/gin"
)

func main() {
	r := server.NewEngine(dummyFetchAssignments)

	r.GET("/taskrader", func(c *gin.Context) {
		c.File("../../../assets/index.html")
	})
	r.GET("/file/main.js", func(c *gin.Context) {
		c.File("../../../assets/main.js")
	})

	host := "localhost"
	port := ":8777"
	r.Run(host + port)
}

func addDateHourMinune(date time.Time, day, hour, minue int) time.Time {
	return date.Add(time.Duration(24*day+hour)*time.Hour + time.Duration(minue)*time.Minute)
}

func dummyFetchAssignments(auth *cred.Credential) ([]*assignment.Assignment, []*assignment.Error) {
	time.Sleep(time.Second * time.Duration(rand.Intn(3)))

	ass := []*assignment.Assignment{}
	errs := []*assignment.Error{}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	if rand.Intn(100) < 60 {
		ass = []*assignment.Assignment{
			{
				Origin:   assignment.OrigGakujo,
				Title:    "第１４回「医学と人間」(2022.1.28)小テスト",
				Course:   "医学と人間",
				Deadline: addDateHourMinune(today, 3, 11, 55), // AM 11:55
			},
			{
				Origin:   assignment.OrigGakujo,
				Title:    "狩野担当分レポート",
				Course:   "知的情報システム開発",
				Deadline: addDateHourMinune(today, 4, 24, 01),
			},
			{
				Origin:   assignment.OrigGakujo,
				Title:    "最終試験レポート",
				Course:   "コンパイラ",
				Deadline: addDateHourMinune(today, 5, 24, 00),
			},
			{
				Origin:   assignment.OrigTeams,
				Title:    "小レポート99",
				Course:   "データベースシステム論",
				Deadline: assignment.UnknownDeadline(),
			},
			{
				Origin:   assignment.OrigEdStem,
				Title:    "†卍小レポート卍†",
				Course:   "人生",
				Deadline: assignment.UnknownDeadline(),
			},
			{
				Origin:   assignment.OrigTeams,
				Title:    "当日課題rp99必須課題",
				Course:   "応用プログラミングC",
				Deadline: addDateHourMinune(today, 0, 23, 55),
			},
			{
				Origin:   assignment.OrigTeams,
				Title:    "当日課題rp99必須課題",
				Course:   "応用プログラミングC",
				Deadline: addDateHourMinune(today, 0, 40, 55),
			},
			{
				Origin:   assignment.OrigTeams,
				Title:    "当日課題rp99必須課題",
				Course:   "応用プログラミングC",
				Deadline: addDateHourMinune(today, 0, -5, 0),
			},
			{
				Origin:   assignment.OrigTeams,
				Title:    "当日課題rp99必須課題",
				Course:   "応用プログラミングC",
				Deadline: now.Add(50 * time.Minute),
			},
			{
				Origin:   assignment.OrigTeams,
				Title:    "当日課題rp99必須課題",
				Course:   "応用プログラミングC",
				Deadline: now.Add(24 * time.Hour),
			},
			{
				Origin:   assignment.OrigTeams,
				Title:    "当日課題rp99必須課題",
				Course:   "応用プログラミングC",
				Deadline: now.Add(48 * time.Hour),
			},
		}
	}
	if rand.Intn(100) < 60 {
		errs = []*assignment.Error{
			{
				Origin: assignment.OrigGakujo,
				Err:    fmt.Errorf("password が空です"),
			},
			{
				Origin: assignment.OrigEdStem,
				Err:    fmt.Errorf("password が空です"),
			},
			{
				Origin: assignment.OrigTeams,
				Err:    fmt.Errorf("password が空です"),
			},
		}
	}

	return ass, errs
}
