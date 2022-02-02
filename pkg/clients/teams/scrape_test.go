package teams

import (
	"testing"
	"time"
)

func tomorrowAt(hour, min int) time.Time {
	t := time.Now().AddDate(0, 0, 1)
	return time.Date(t.Year(), t.Month(), t.Day(), hour, min, 0, 0, time.Local)
}

func todayAt(hour, min int) time.Time {
	t := time.Now()
	return time.Date(t.Year(), t.Month(), t.Day(), hour, min, 0, 0, time.Local)
}

func yesterdayAt(hour, min int) time.Time {
	t := time.Now().AddDate(0, 0, -1)
	return time.Date(t.Year(), t.Month(), t.Day(), hour, min, 0, 0, time.Local)
}

func TestParseDueText(t *testing.T) {
	testcases := []struct {
		Text     string
		Expected time.Time
	}{
		{
			Text:     "Due February 1, 2022 11:59 PM",
			Expected: time.Date(2022, time.February, 1, 23, 59, 0, 0, time.Local),
		},
		{
			Text:     "Due tomorrow at 2:59 PM",
			Expected: tomorrowAt(14, 59),
		},
		{
			Text:     "Due today at 3:00 PM",
			Expected: todayAt(15, 0),
		},
		{
			Text:     "Due yesterday at 11:59 PM",
			Expected: yesterdayAt(23, 59),
		},
		{
			Text:     "期限 2022年1月26日 23:59",
			Expected: time.Date(2022, time.January, 26, 23, 59, 0, 0, time.Local),
		},
		{
			Text:     "明日 00:00 が期限",
			Expected: tomorrowAt(0, 0),
		},
		{
			Text:     "今日 23:00 が期限",
			Expected: todayAt(23, 0),
		},
		{
			Text:     "昨日 23:59 が期限",
			Expected: yesterdayAt(23, 59),
		},
	}

	for _, c := range testcases {
		got, err := parseDueText(c.Text)
		if err != nil {
			t.Error(err)
		} else if !got.Equal(c.Expected) {
			t.Errorf("\n  Input string: %q\n  got : %s\n  want: %s", c.Text, got, c.Expected)
		}
	}
}
