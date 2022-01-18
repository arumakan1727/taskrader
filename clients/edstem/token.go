package edstem

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) getToken(email, password string) (string, error) {
	jsonclient := JsonClient{Email: email, Password: password}

	jsonjson, _ := json.Marshal(jsonclient)

	resp, err := c.client.Post("https://edstem.org/api/token", "application/json", bytes.NewBuffer(jsonjson))
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Response status was %d(except %d)", resp.StatusCode, http.StatusOK)
	}

	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	token := Token{}
	err = json.Unmarshal(res, &token)
	if err != nil {
		return "", err
	}

	return token.Token, nil
}
