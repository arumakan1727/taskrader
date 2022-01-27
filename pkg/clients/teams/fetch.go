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
	selectorUseWebAppLink    = `#download-desktop-page a.use-app-lnk`
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
	if err := driver.Start(); err != nil {
		return "", err
	}
	defer driver.Stop()

	page, err := driver.NewPage()
	if err != nil {
		return "", err
	}

	logger.Printf("Opening %s ...", teamsURL)
	if err := page.Navigate(teamsURL); err != nil {
		return "", err
	}
	time.Sleep(3 * time.Second)

	{
		url, err := page.URL()
		if err != nil {
			return "", err
		}
		if strings.HasPrefix(url, loginURL) {
			return "", NewErrLoginRequired()
		}
	}

	// たまにデスクトップ版を宣伝する画面が出てくるので、その場合はWeb版続行リンクを押下して数秒待つ
	if err := clickElemBySelector(page, selectorUseWebAppLink, 2*time.Second); err == nil {
		time.Sleep(4 * time.Second)
	}

	logger.Println("Switching to assignments tab ...")
	if err := clickElemBySelector(page, selectorAssignmentTabBtn, 10*time.Second); err != nil {
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
		return nil, fmt.Errorf("failed to find assignment-list: `%s`", selectorAssignmentList)
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
		logger.Printf("dueText = %q, err=%v, parseResult=%q", dueText, err, tm)

		res = append(res, Assignment{
			Title:    title,
			Course:   course,
			Deadline: tm,
		})
	})

	return res, nil
}

func parseDueText(text string) (time.Time, error) {
	if strings.HasPrefix(text, "Due") {
		if strings.HasPrefix(text, "Due tomorrow at") {
			replaceTo := time.Now().AddDate(0, 0, 1).Format("January 2, 2006")
			text = strings.Replace(text, "tomorrow at", replaceTo, 1)

		} else if strings.HasPrefix(text, "Due today at") {
			replaceTo := time.Now().Format("January 2, 2006")
			text = strings.Replace(text, "today at", replaceTo, 1)
		}
		return time.ParseInLocation("Due January 2, 2006 3:04 PM", text, time.Local)
	}
	if strings.HasPrefix(text, "明日") {
		text = strings.TrimPrefix(text, "明日")
		text = strings.TrimSuffix(text, "が期限")
		text = strings.TrimSpace(text)
		t := time.Now().AddDate(0, 0, 1)
		text = "期限" + t.Format(" 2006年1月2日 ") + text
	}
	if strings.HasPrefix(text, "今日") {
		text = strings.TrimPrefix(text, "今日")
		text = strings.TrimSuffix(text, "が期限")
		text = strings.TrimSpace(text)
		t := time.Now()
		text = "期限" + t.Format(" 2006年1月2日 ") + text
	}
	if strings.HasPrefix(text, "期限") {
		return time.ParseInLocation("期限 2006年1月2日 15:04", text, time.Local)
	}
	return time.Time{}, fmt.Errorf("unknown dueText format: %s", text)
}
