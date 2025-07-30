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
	client    *http.Client
	urlSchema string
	host      string
	apiToken  string
}

type accountRequest struct {
	Username     string       `json:"userName"`
	Password     string       `json:"password"`
	RobotEdition RobotEdition `json:"robotEdition,omitempty"`
}

type AccountResponse struct {
	Token string `json:"token"`
}

type projectInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func NewBotApiClient(host string, urlSchema string, username string, password string, robotEdition RobotEdition) (*BotApiClient, error) {
	botApiClient := new(BotApiClient)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c := &http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}

	botApiClient.client = c
	botApiClient.host = host
	botApiClient.urlSchema = urlSchema

	token, err := botApiClient.postAccount(username, password, robotEdition)
	if err != nil {
		return botApiClient, err
	}
	botApiClient.apiToken = token

	return botApiClient, nil
}

func (b *BotApiClient) postAccount(username string, password string, robotEdition RobotEdition) (string, error) {
	var token string

	url := url.URL{
		Scheme: b.urlSchema,
		Host:   b.host,
		Path:   path.Join("api", "Account"),
	}
	a := accountRequest{
		Username:     username,
		Password:     password,
		RobotEdition: robotEdition,
	}

	accountJson, err := json.Marshal(a)
	if err != nil {
		return token, err
	}

	bodyBytes := bytes.NewReader(accountJson)
	response, err := b.client.Post(url.String(), "application/json", bodyBytes)
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
	token = respBody.Token

	return token, nil
}

func (b *BotApiClient) GetProjects() ([]projectInfo, error) {
	projects := make([]projectInfo, 0)
	url := &url.URL{
		Scheme: b.urlSchema,
		Host:   b.host,
		Path:   path.Join("api", "RpaProjects", "v2"),
	}

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return projects, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", b.apiToken))
	req.Header.Set("Accept", "text/plain")

	response, err := b.client.Do(req)
	if err != nil {
		return projects, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return projects, fmt.Errorf("'%v: %v' request error. Code: %v", req.Method, url.Path, response.StatusCode)
	}

	respBodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return projects, err
	}

	err = json.Unmarshal(respBodyBytes, &projects)
	if err != nil {
		return projects, err
	}

	return projects, nil
}
