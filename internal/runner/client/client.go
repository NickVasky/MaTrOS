package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"
)

type RobotEdition uint

const (
	Enterprise RobotEdition = 1 + iota
	Standard
	Desktop
)

type BotApiClient struct {
	Client    *http.Client
	URLSchema string
	Host      string
}

type AccountRequest struct {
	Username     string       `json:"userName"`
	Password     string       `json:"password"`
	RobotEdition RobotEdition `json:"robotEdition,omitempty"`
}

type AccountResponse struct {
	Token string `json:"token"`
}

func NewBotApiClient(Host string, URLSchema string) *BotApiClient {
	botApiClient := new(BotApiClient)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c := &http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}

	botApiClient.Client = c
	botApiClient.Host = Host
	botApiClient.URLSchema = URLSchema

	return botApiClient
}

func (b *BotApiClient) PostAccount(user string, pass string, robotEdition RobotEdition) (string, error) {
	var token string

	u := url.URL{
		Scheme: b.URLSchema,
		Host:   b.Host,
		Path:   path.Join("api", "Account"),
	}
	a := AccountRequest{
		Username:     user,
		Password:     pass,
		RobotEdition: robotEdition,
	}

	accountJson, err := json.Marshal(a)
	if err != nil {
		return token, err
	}

	bodyBytes := bytes.NewReader(accountJson)
	response, err := b.Client.Post(u.String(), "application/json", bodyBytes)
	if err != nil {
		return token, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return token, fmt.Errorf("Authtorization error. Code: %v", response.StatusCode)
	}

	respBodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return token, err
	}
	var respBody AccountResponse
	err = json.Unmarshal(respBodyBytes, &respBody)
	if err != nil {
		return token, err
	}
	token = fmt.Sprintf("Bearer: %v", respBody.Token)

	return token, nil
}
