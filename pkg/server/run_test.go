package server_test

import (
	"os"
	"testing"

	"github.com/arumakan1727/taskrader/pkg/server"
)

func TestRunServer(t *testing.T) {
	if os.Getenv("NOW_ON_CI") != "" {
		return
	}

	r := server.NewEngine()

	host := "localhost"
	port := ":8777"
	r.Run(host + port)
}
