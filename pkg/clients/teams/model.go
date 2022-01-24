package teams

import "time"

type Assignment struct {
	Title    string    // 課題タイトル
	Course   string    // 課題の科目
	Deadline time.Time // 締切
}
