package teams

import (
	"fmt"
	"time"

	"github.com/sclevine/agouti"
)

func getSelectionValue(elem *agouti.Selection) string {
	value, err := elem.Attribute("value")
	if err != nil || value == "" {
		return ""
	}
	return value
}

func sendKeys(page *agouti.Page, name, keys string, timeout time.Duration) error {
	startTime := time.Now()
	for getSelectionValue(page.FindByName(name)) == "" {
		if time.Since(startTime) > timeout {
			return fmt.Errorf("Failed to send keys to element having name=%s", name)
		}
		el := page.FindByName(name)
		el.Clear()
		time.Sleep(200 * time.Millisecond)
		el.SendKeys(keys)
		time.Sleep(200 * time.Millisecond)
	}
	return nil
}

func clickButtonHavingID(page *agouti.Page, buttonID string, timeout time.Duration) error {
	startTime := time.Now()
	for {
		if err := page.FindByID(buttonID).Click(); err == nil {
			return nil
		}
		time.Sleep(200 * time.Millisecond)
		if time.Since(startTime) > timeout {
			return fmt.Errorf("Failed to click button having id=%s", buttonID)
		}
	}
}

func existsElementHavingID(page *agouti.Page, id string, interval, timeout time.Duration) bool {
	startTime := time.Now()
	for {
		count, err := page.FindByID(id).Count()
		if err == nil {
			return count > 0
		}
		if time.Since(startTime) > timeout {
			return false
		}
		time.Sleep(interval)
	}
}
