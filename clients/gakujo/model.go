package gakujo

import "time"

type TaskRow struct {
	Type     string
	Deadline time.Time
	Name     string
}
