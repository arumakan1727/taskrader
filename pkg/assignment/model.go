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
	Origin   Origin    // 課題の出所
	Title    string    // 課題のタイトル
	Course   string    // 科目orコース名
	Deadline time.Time // 課題の締め切り; 不明の場合は UnknownDeadline() を設定する
}

var _unknownDeadline = time.Unix(1<<60, 999999999)

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
