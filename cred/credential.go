package cred

import (
	"os"
)

type Credential struct {
	Gakujo Gakujo
	EdStem EdStem
	Teams  Teams
}

type Gakujo struct {
	Username string
	Password string
}

type EdStem struct {
	Email    string
	Password string
}

type Teams struct {
	Token string
}

// 環境変数から Credential を生成する。
// .env は自動では読み込まないので適宜 godotenv.Load() すること。
// 該当する名前の環境変数が存在しない場合は、そのメンバ変数の値は空文字列となる。
func LoadFromEnv() *Credential {
	return &Credential{
		Gakujo: Gakujo{
			Username: os.Getenv("GAKUJO_USERNAME"),
			Password: os.Getenv("GAKUJO_PASSWORD"),
		},
		EdStem: EdStem{
			Email:    os.Getenv("EDSTEM_EMAIL"),
			Password: os.Getenv("EDSTEM_PASSWORD"),
		},
		Teams:  Teams{
			Token: os.Getenv("TEAMS_TOKEN"),
		},
	}
}
