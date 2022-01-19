package gakujo

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func (c *Client) getTask() []TaskRow {

	datas := make(url.Values)
	datas.Set("headTitle", "ホーム")
	datas.Set("menuCode", "Z07") // TODO: menucode を定数化(まとめてやる)
	datas.Set("nextPath", "/home/home/initialize")

	urll := "https://gakujo.shizuoka.ac.jp/portal/common/generalPurpose/"

	resp, err := c.getPage(urll, datas)

	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(resp)

	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(io.NopCloser(bytes.NewBuffer(body)))

	if err != nil {
		log.Fatal(err)
	}

	var taskRows []TaskRow

	//taskRows = make([]TaskRow, 0)
	doc.Find("#tbl_submission > tbody > tr").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		var inerr error
		taskType := selection.Find("td.arart > span > span").Text()
		deadlineText := selection.Find("td.daytime").Text()
		var deadline time.Time
		if deadlineText != "" {
			deadline, inerr = Parse2400("2006/01/02 15:04", deadlineText)
			if inerr != nil {
				err = inerr
				return false
			}
		}
		data := TaskRow{
			Type:     taskType,
			Deadline: deadline,
			Name:     selection.Find("td:nth-child(3) > a").Text(),
		}
		taskRows = append(taskRows, data)
		return true
	})

	return taskRows
}

func Parse2400(layout, value string) (time.Time, error) {
	parsedTime, err := time.Parse(layout, value)
	if err != nil {
		if !isHourOutErr(err) {
			return time.Time{}, err
		}
		i := strings.Index(layout, "15")
		if i == -1 {
			return time.Time{}, errors.New("stdHour 15 was not found in layout")
		}
		newValue := value[:i] + "00" + value[i+2:]
		parsedTime, err = time.Parse(layout, newValue)
		if err != nil {
			return time.Time{}, err
		}
		return parsedTime.Add(24 * time.Hour), nil
	}
	return parsedTime, nil
}

func isHourOutErr(err error) bool {
	switch err.(type) {
	case *time.ParseError:
		return strings.Contains(err.Error(), "hour")
	default:
		return false
	}
}
