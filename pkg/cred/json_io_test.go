package cred_test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/arumakan1727/taskrader/pkg/cred"
)

var jsonPath string

func TestMain(m *testing.M) {
	// 一時的なディレクトリ & json ファイルパスの取得
	tempDir, err := ioutil.TempDir(os.TempDir(), "gotest")
	if err != nil {
		panic(err)
	}
	jsonPath = path.Join(tempDir, "hoge", "taskrader", "credential.json")

	// テスト実行
	m.Run()

	// 一時的に生成したディレクトリを再帰的に削除
	os.RemoveAll(tempDir)
}

func TestSaveAndLoadJSON(t *testing.T) {
	c := cred.Credential{
		Gakujo: cred.Gakujo{
			Username: "gakujo-username",
			Password: "gakujo-password",
		},
		EdStem: cred.EdStem{
			Email:    "edstem_email@example.com",
			Password: "edstem-password",
		},
		Teams: cred.Teams{
			Email:    "teams_email@example.com",
			Password: "teams-password",
		},
	}

	t.Logf("Temporary jsonPath = %s", jsonPath)

	t.Run("SaveToJSONFile() should success", func(t *testing.T) {
		if err := c.SaveToJSONFile(jsonPath); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("LoadFromJSONFileOrEmpty() should return credential which is equivalent to saved value", func(t *testing.T) {
		got, err := cred.LoadFromJSONFile(jsonPath)
		if err != nil {
			t.Error(err)
		}

		if got.Gakujo.Username != c.Gakujo.Username {
			t.Errorf("Expected %q, but got %q", c.Gakujo.Username, got.Gakujo.Username)
		}
		if got.Gakujo.Password != c.Gakujo.Password {
			t.Errorf("Expected %q, but got %q", c.Gakujo.Password, got.Gakujo.Password)
		}

		if got.EdStem.Email != c.EdStem.Email {
			t.Errorf("Expected %q, but got %q", c.EdStem.Email, got.EdStem.Email)
		}
		if got.EdStem.Password != c.EdStem.Password {
			t.Errorf("Expected %q, but got %q", c.EdStem.Password, got.EdStem.Password)
		}

		if got.Teams.Email != c.Teams.Email {
			t.Errorf("Expected %q, but got %q", c.Teams.Email, got.Teams.Email)
		}
		if got.Teams.Password != c.Teams.Password {
			t.Errorf("Expected %q, but got %q", c.Teams.Password, got.Teams.Password)
		}
	})
}
