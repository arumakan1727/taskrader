package cred

import (
	"fmt"
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
	Email    string
	Password string
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
		Teams: Teams{
			Email:    os.Getenv("TEAMS_EMAIL"),
			Password: os.Getenv("TEAMS_PASSWORD"),
		},
	}
}

type ErrEmpty struct {
	FieldName string
}

func (e *ErrEmpty) Error() string {
	return fmt.Sprintf("%s is empty", e.FieldName)
}

func newErrEmpty(fieldName string) ErrEmpty {
	return ErrEmpty{
		FieldName: fieldName,
	}
}

func (c *Credential) CheckEmptyField() []ErrEmpty {
	errs := make([]ErrEmpty, 0)
	if c.Gakujo.Username == "" {
		errs = append(errs, newErrEmpty("Gakujo.Username"))
	}
	if c.Gakujo.Password == "" {
		errs = append(errs, newErrEmpty("Gakujo.Username"))
	}
	if c.EdStem.Email == "" {
		errs = append(errs, newErrEmpty("EdStem.Email"))
	}
	if c.EdStem.Password == "" {
		errs = append(errs, newErrEmpty("EdStem.Password"))
	}
	if c.Teams.Email == "" {
		errs = append(errs, newErrEmpty("Teams.Email"))
	}
	if c.Teams.Password == "" {
		errs = append(errs, newErrEmpty("Teams.Password"))
	}
	return errs
}

func (c *Credential) AbortIfEmptyFieldExists() {
	if errs := c.CheckEmptyField(); len(errs) > 0 {
		for _, e := range errs {
			fmt.Fprintln(os.Stderr, "Error: "+e.Error())
		}
		os.Exit(1)
	}
}
