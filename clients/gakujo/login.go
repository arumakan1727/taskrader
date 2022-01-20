package gakujo

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	HostName          = "https://gakujo.shizuoka.ac.jp"
	IdpHostName       = "https://idp.shizuoka.ac.jp"
	GeneralPurposeUrl = "https://gakujo.shizuoka.ac.jp/portal/common/generalPurpose/"
)

type Client struct {
	client *http.Client
	jar    *cookiejar.Jar
	token  string // org.apache.struts.taglib.html.TOKEN
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
		Timeout: 5 * time.Minute,
	}
	return &Client{
		client: &httpClient,
		jar:    jar,
	}
}

func (c *Client) Login(username, password string) error {
	if err := c.fetchGakujoPortalJSESSIONID(); err != nil {
		return err
	}

	if err := c.fetchGakujoRootJSESSIONID(); err != nil {
		return err
	}

	if err := c.preLogin(); err != nil {
		return err
	}
	resp, err := c.shibbolethlogin()
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusFound {
		return fmt.Errorf("Response status was %d(expect %d or %d)", resp.StatusCode, http.StatusOK, http.StatusFound)
	}

	// セッションがないとき
	if resp.StatusCode == http.StatusFound {
		loginAPIurl, err := c.fetchLoginAPIurl(resp.Header.Get("Location"))
		if err != nil {
			return err
		}
		if err := c.login(IdpHostName+loginAPIurl, username, password); err != nil {
			return err
		}
	}

	return c.initialize()
}

func (c *Client) fetchGakujoPortalJSESSIONID() error {
	resp, err := c.get("https://gakujo.shizuoka.ac.jp/portal/")
	if err != nil {
		return err
	}
	defer func() {
		resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
	}()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Response status was %d(expect %d)", resp.StatusCode, http.StatusOK)
	}

	return nil
}

func (c *Client) fetchGakujoRootJSESSIONID() error {
	unixmilli := time.Now().UnixNano() / 1000000
	resp, err := c.get("https://gakujo.shizuoka.ac.jp/UI/jsp/topPage/topPage.jsp?_=" + strconv.FormatInt(unixmilli, 10))
	if err != nil {
		return err
	}
	defer func() {
		resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
	}()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Response status was %d(expect %d)", resp.StatusCode, http.StatusOK)
	}

	return nil
}

func (c *Client) preLogin() error {
	datas := url.Values{}
	datas.Set("mistakeChecker", "0")

	resp, err := c.client.PostForm("https://gakujo.shizuoka.ac.jp/portal/login/preLogin/preLogin", datas)
	if err != nil {
		return err
	}
	defer func() {
		resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
	}()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Response status was %d(expect %d)", resp.StatusCode, http.StatusOK)
	}

	return nil
}

func (c *Client) fetchLoginAPIurl(SSOSAMLRequestURL string) (string, error) {
	resp, err := c.get(SSOSAMLRequestURL)
	if err != nil {
		return "", err
	}
	defer func() {
		resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
	}()
	if resp.StatusCode != http.StatusFound {
		return "", fmt.Errorf("Response status was %d(expect %d)", resp.StatusCode, http.StatusOK)
	}
	return resp.Header.Get("Location"), nil
}

func RelayStateAndSAMLResponse(htmlReader io.ReadCloser) (string, string, error) {
	doc, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		return "", "", err
	}
	selection := doc.Find("html > body > form > div > input")
	relayState, ok := selection.Attr("value")
	if !ok {
		return "", "", errors.New("RelayState")
	}
	selection = selection.Next()
	samlResponse, ok := selection.Attr("value")
	if !ok {
		return "", "", errors.New("SAMLResponse")
	}

	return relayState, samlResponse, nil
}

func (c *Client) login(reqUrl, username, password string) error {
	htmlReadCloser, err := c.postSSOexecution(reqUrl, username, password)
	if err != nil {
		return err
	}
	relayState, samlResponse, err := RelayStateAndSAMLResponse(htmlReadCloser)
	if err != nil {
		return err
	}
	htmlReadCloser.Close()
	_, _ = io.Copy(io.Discard, htmlReadCloser)

	location, err := c.fetchSSOinitLoginLocation(relayState, samlResponse)
	if err != nil {
		return err
	}

	resp, err := c.getWithReferer(location, "https://idp.shizuoka.ac.jp/")
	if err != nil {
		return err
	}
	defer func() {
		resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
	}()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Response status was %d(expect %d)", resp.StatusCode, http.StatusOK)
	}

	return nil
}

func (c *Client) postSSOexecution(reqUrl, username, password string) (io.ReadCloser, error) {
	datas := make(url.Values)
	datas.Set("j_username", username)
	datas.Set("j_password", password)
	datas.Set("_eventId_proceed", "")

	resp, err := c.client.PostForm(reqUrl, datas)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Response status was %d(expect %d)", resp.StatusCode, http.StatusOK)
	}

	return resp.Body, nil
}

func (c *Client) shibbolethlogin() (*http.Response, error) {
	url := HostName + "/portal/shibbolethlogin/shibbolethLogin/initLogin/sso"
	req, _ := http.NewRequest(http.MethodPost, url, nil)
	resp, err := c.request(req)
	resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	return resp, err
}

func (c *Client) fetchSSOinitLoginLocation(relayState, samlResponse string) (string, error) {
	reqUrl := "https://gakujo.shizuoka.ac.jp/Shibboleth.sso/SAML2/POST"

	datas := make(url.Values)
	datas.Set("RelayState", relayState)
	datas.Set("SAMLResponse", samlResponse)

	resp, err := c.client.PostForm(reqUrl, datas)
	if err != nil {
		return "", err
	}
	defer func() {
		resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
	}()
	if resp.StatusCode != http.StatusFound {
		return "", fmt.Errorf("%s\nResponse status was %d(expect %d)", reqUrl, resp.StatusCode, http.StatusFound)
	}

	return resp.Header.Get("Location"), nil
}

func (c *Client) initialize() error {
	reqURL := "https://gakujo.shizuoka.ac.jp/portal/home/home/initialize"

	datas := make(url.Values)
	datas.Set("EXCLUDE_SET", "")

	rc, err := c.getPage(reqURL, datas)
	if err != nil {
		return err
	}
	defer func() {
		rc.Close()
		_, _ = io.Copy(io.Discard, rc)
	}()

	return nil
}

// search a cookie "JSESSIONID" from c.jar
// if not found, return ""
func (c *Client) SessionID() string {
	u, _ := url.Parse(HostName)
	for _, cookie := range c.jar.Cookies(u) {
		if cookie.Name == "JSESSIONID" {
			return cookie.Value
		}
	}

	return ""
}

// save cookie "Set-Cookies" into client.cookie
func (c *Client) request(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func ApacheToken(htmlReader io.ReadCloser) (string, error) {
	// ページによってtokenの場所が違う場合
	selectors := []string{
		"#SC_A01_06 > form:nth-child(15) > div > input[type=hidden]",
		"#header > form:nth-child(4) > div > input[type=hidden]",
	}
	doc, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		return "", err
	}
	for _, selector := range selectors {
		selection := doc.Find(selector)
		token, ok := selection.Attr("value")
		if ok {
			return token, nil
		}
	}
	return "", errors.New("org.apache.struts.taglib.html.TOKEN")
}

// get page which needs org.apache.struts.taglib.html.TOKEN and save its token
func (c *Client) getPage(url string, datas url.Values) (io.ReadCloser, error) {
	datas.Set("org.apache.struts.taglib.html.TOKEN", c.token)

	resp, err := c.postForm(url, datas)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Response status was %d(expext %d)", resp.StatusCode, http.StatusOK)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	token, err := ApacheToken(io.NopCloser(bytes.NewReader(b)))
	if err != nil {
		// getPage では必ず apache Token が含まれるページを取得するはず
		return nil, err
	}
	c.token = token

	return io.NopCloser(bytes.NewReader(b)), nil
}

// http.Get wrapper
func (c *Client) get(url string) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	return c.request(req)
}

// http.Get wrapper
func (c *Client) getWithReferer(url, referer string) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Referer", referer)
	return c.request(req)
}

// http.PostForm wrapper
func (c *Client) postForm(url string, datas url.Values) (*http.Response, error) {
	req, err := http.NewRequest(
		http.MethodPost,
		url,
		strings.NewReader(datas.Encode()),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.request(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
