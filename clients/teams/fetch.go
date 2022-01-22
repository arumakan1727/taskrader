package teams

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sclevine/agouti"
)

const (
	selectorTeamTabBtn       = `nav > ul > li:nth-child(3) > button`
	selectorAssignmentTabBtn = `nav > ul > li:nth-child(4) > button`
	selectorIFrame           = `body > app-caching-container > div > div > extension-tab > div > embedded-page-container > div > iframe`
	selectorAssignmentList   = `#root div.ms-FocusZone[data-test="assignment-list"]`
	selectorAssignmentCard   = `div[data-test="assignment-card"]`
	selectorCardTitle        = `div > div > div > div:nth-child(2)`
)

func FetchAssignments(logger *log.Logger) ([]Assignment, error) {
	html, err := fetchAssignmentsPageHTML(logger)
	if err != nil {
		return nil, err
	}
	return scrapeAssignmentList(html, logger)
}

func fetchAssignmentsPageHTML(logger *log.Logger) (string, error) {
	opt, err := myChromeOptions()
	if err != nil {
		return "", err
	}
	driver := agouti.ChromeDriver(opt)
	defer driver.Stop()

	if err := driver.Start(); err != nil {
		return "", err
	}
	page, err := driver.NewPage()
	if err != nil {
		return "", err
	}

	logger.Printf("Opening %s ...", teamsURL)
	if err := page.Navigate(teamsURL); err != nil {
		return "", err
	}
	time.Sleep(2 * time.Second)

	{
		url, err := page.URL()
		if err != nil {
			return "", err
		}
		if strings.HasPrefix(url, loginURL) {
			return "", NewErrLoginRequired()
		}
	}

	logger.Println("Switching to assignments tab ...")
	if err := clickButtonBySelector(page, selectorAssignmentTabBtn, 3*time.Second); err != nil {
		return "", err
	}
	time.Sleep(3 * time.Second)
	if err := page.First(selectorIFrame).SwitchToFrame(); err != nil {
		return "", err
	}

	logger.Println("Waiting assignment list ...")
	startTime := time.Now()
	timeout := 15 * time.Second
	for {
		cards := page.First(selectorAssignmentList).All(selectorAssignmentCard)
		count, err := cards.Count()
		logger.Println("cards.Count() = ", count, ", err = ", err)
		if err == nil {
			break
		}
		if time.Since(startTime) > timeout {
			return "", fmt.Errorf("assignment-card cannot be found")
		}
		time.Sleep(time.Second)
	}

	return page.HTML()
}

func scrapeAssignmentList(html string, logger *log.Logger) ([]Assignment, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	listContainer := doc.Find(selectorAssignmentList)
	if attr, exists := listContainer.Attr("data-test"); !exists || attr != "assignment-list" {
		return nil, fmt.Errorf("Failed to find assignment-list: `%s`", selectorAssignmentList)
	}

	cardElems := listContainer.Find(`a > div[data-test="assignment-card"]`)
	log.Printf("cardElemes.Length() = %d\n", cardElems.Length())

	res := make([]Assignment, 0, cardElems.Length())
	listContainer.Find(`a > div[data-test="assignment-card"]`).Each(func(i int, card *goquery.Selection) {
		el := card.Find(selectorCardTitle)
		title := strings.TrimSpace(el.Text())

		el = el.Next()
		course := strings.TrimSpace(el.Text())

		el = el.Next()
		dueText := strings.TrimSpace(el.Find(`span > span:nth-child(3)`).Text())
		tm, err := parseDueText(dueText)
		if err != nil {
			tm = time.Time{}
		}
		logger.Printf("dueText = %q, err=%v, tm=%q", dueText, err, tm)

		res = append(res, Assignment{
			Title:    title,
			Course:   course,
			Deadline: tm,
		})
	})

	return res, nil
}

func parseDueText(text string) (time.Time, error) {
	return time.ParseInLocation("Due January 2, 2006 03:04 PM", text, time.Local)
}
