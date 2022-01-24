package config

import (
	"os"
	"path"
)

func TaskRaderCacheDir() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return path.Join(cacheDir, "taskrader"), nil
}
