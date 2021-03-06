package teams

import (
	"os"
	"path"

	"github.com/arumakan1727/taskrader/pkg/config"
	"github.com/sclevine/agouti"
)

const (
	loginURL = `https://login.microsoftonline.com`
	teamsURL = "https://teams.microsoft.com"
)

// Selenium の ChromeDriver が使うユーザプロフィールデータの保存先ディレクトリ
// Cookie はここに永続保存される
func ChromeTmpUserDataDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dirInfo, err := os.Stat(path.Join(homeDir, "snap", "chromium"))
	if err == nil && dirInfo.IsDir() {
		return path.Join(homeDir, "snap", "chromium", "current", ".cache", "taskrader", "chrome-tmp-profile"), nil
	}

	cacheDir, err := config.TaskraderCacheDir()
	if err != nil {
		return "", err
	}
	return path.Join(cacheDir, "chrome-tmp-profile"), nil
}

func ClearCookies() {
	dir, err := ChromeTmpUserDataDir()
	if err != nil {
		return
	}
	os.Remove(path.Join(dir, "Default", "Cookies"))
	os.Remove(path.Join(dir, "Default", "Cookies-journal"))
	os.RemoveAll(path.Join(dir, "Default", "Sessions"))
	os.RemoveAll(path.Join(dir, "Default", "Session Storage"))
}

func myChromeOptions() (agouti.Option, error) {
	dir, err := ChromeTmpUserDataDir()
	if err != nil {
		return nil, err
	}

	opt := agouti.ChromeOptions(
		"args", []string{
			"--headless",
			"--disable-gpu",
			"--user-data-dir=" + dir,
			"--lang=en",
		},
	)
	return opt, nil
}
