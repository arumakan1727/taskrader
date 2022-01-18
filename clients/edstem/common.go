package edstem

import (
	"io"
	"net/http"
)

func (c *Client) GetJson(method, url string) ([]byte, error) {
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Set("x-token", c.token)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return res, nil
}
