package gakujo

import "time"

type TaskRow struct {
	Type     string    // レポート, 小テスト, etc
	Deadline time.Time // 締切
	Title    string    // 課題のタイトル
	Course   string    // 課題の科目名
}
