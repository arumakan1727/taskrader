package config_test

import (
	"os"
	"path"
	"testing"

	"github.com/arumakan1727/taskrader/pkg/config"
)

func TestTaskraderCacheDir(t *testing.T) {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		t.Fatal(err)
	}

	want := path.Join(userCacheDir, "taskrader")
	got, err := config.TaskraderCacheDir()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("taskraderCacheDir = %s", got)

	if got != want {
		t.Errorf("Expected cache dir = %s, but got %s", want, got)
	}
}

func TestTaskraderConfigDir(t *testing.T) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		t.Fatal(err)
	}

	want := path.Join(userConfigDir, "taskrader")
	got, err := config.TaskraderConfigDir()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("taskraderConfigDir = %s", got)

	if got != want {
		t.Errorf("Expected config dir = %s, but got %s", want, got)
	}
}
