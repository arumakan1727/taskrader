package edstem

import (
	"net/http"
	"net/http/cookiejar"
	"time"
)

type Client struct {
	client *http.Client
	jar    *cookiejar.Jar
	token  string
}

func NewClient() *Client {
	jar, _ := cookiejar.New(
		nil,
	)
	httpClient := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar:     jar,
		Timeout: time.Second * 3,
	}
	return &Client{
		client: &httpClient,
		jar:    jar,
	}
}

func (c *Client) Login(email, password string) error {
	var err error
	c.token, err = c.getToken(email, password)
	if err != nil {
		return err
	}
	return nil
}