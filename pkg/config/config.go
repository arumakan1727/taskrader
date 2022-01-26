package config

import (
	"os"
	"path"
)

func TaskraderCacheDir() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return path.Join(cacheDir, "taskrader"), nil
}

func TaskraderConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return path.Join(configDir, "taskrader"), nil
}

func TaskraderCredentialJSONPath() (string, error) {
	configDir, err := TaskraderConfigDir()
	if err != nil {
		return "", err
	}
	return path.Join(configDir, "credential.json"), nil
}
