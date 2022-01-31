package assignment

import (
	"fmt"
	"time"
)

// 課題の出所を表す enum のような型
type Origin string

const (
	OrigGakujo = Origin("学情")
	OrigTeams  = Origin("Teams")
	OrigEdStem = Origin("EdStem")
)

// 課題を表す共通の型
type Assignment struct {
	Origin   Origin    `json:"origin"` // 課題の出所
	Title    string    `json:"title"`  // 課題のタイトル
	Course   string    `json:"course"` // 科目orコース名
	Deadline time.Time `json:"due"`    // 課題の締め切り; 不明の場合は UnknownDeadline() を設定する
}

var _unknownDeadline = time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)

// 不明な締切日時を表す定数
func UnknownDeadline() time.Time {
	return _unknownDeadline
}

// デバッグで Println(a) できるようにするために定義しておく; ユーザへの表示の際は使わない
func (a *Assignment) String() string {
	deadline := "不明"
	if a.Deadline != UnknownDeadline() {
		deadline = a.Deadline.String()
	}

	return fmt.Sprintf("Assignment{Origin=%s, Title=%q, Course=%q, deadline=%q}", a.Origin, a.Title, a.Course, deadline)
}

type Error struct {
	Origin Origin
	Err    error
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Origin, e.Err.Error())
}
