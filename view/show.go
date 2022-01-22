package view

import (
	"fmt"

	"sort"
	"time"

	"github.com/arumakan1727/taskrader/assignment"
	color "github.com/fatih/color"
)

func Show(tasks []*assignment.Assignment) {
	const formatHour = "2006/01/02"
	const formatMin = "15:04"
	wdays := []string{"日", "月", "火", "水", "木", "金", "土"}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Deadline.Before(tasks[j].Deadline)
	})

	line := makeLine(tasks)

	fmt.Println(len(tasks), "件の未提出課題：")
	for i := 0; i < len(tasks); i++ {

		limit := "<締切不明>"

		if tasks[i].Deadline != time.Unix(1<<60, 999999999) {
			limit = tasks[i].Deadline.Format(formatHour) + " (" + wdays[tasks[i].Deadline.Weekday()] + ") " + tasks[i].Deadline.Format(formatMin)
		}

		fmt.Println(line)
		output := fmt.Sprintf("| [%d] %s |\n|     %s @%s |", i+1, tasks[i].Title, limit, tasks[i].Origin)
		printWithColor(output, tasks[i].Deadline)
	}

	fmt.Println(line)

}

func makeLine(tasks []*assignment.Assignment) string {
	max := 0

	for i := 0; i < len(tasks); i++ {
		if max < len(tasks[i].Title) {
			max = len(tasks[i].Title)
		}
	}

	line := ""

	for i := 0; i < max+4; i++ {
		line += "-"
	}

	return line
}

func printWithColor(output string, deadline time.Time) {
	if diff := time.Until(deadline).Hours(); 0 <= diff && diff <= 24 {
		fmt.Println(color.RedString(output))
	} else if 24 < diff && diff <= 48 {
		fmt.Println(color.YellowString(output))
	} else {
		fmt.Println(output)
	}
}
