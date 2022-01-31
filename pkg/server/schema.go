package server

import "github.com/arumakan1727/taskrader/pkg/assignment"

type RespAssErr struct {
	Origin  string `json:"origin"`
	Message string `json:"message"`
}

type RespAssignmentsAndErrors struct {
	Ass    []*assignment.Assignment `json:"assignments"`
	Errors []RespAssErr             `json:"errors"`
}

type RespSimpleErr struct {
	Message string `json:"errmsg"`
}

type RespLoginStatus struct {
	GakujoLogined bool `json:"gakujo"`
	EdstemLogined bool `json:"edstem"`
	TeamsLogined  bool `json:"teams"`
}
