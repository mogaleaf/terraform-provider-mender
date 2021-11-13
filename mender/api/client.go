package api

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
)

var (
	loginUrl  = "api/management/v1/useradm/auth/login"
	uploadUrl = "api/management/v1/deployments/artifacts"
)

type Client interface {
	Login() error
	UploadArtifact([]byte) (string, error)
}

func New(host, username, password string) Client {
	return &client{
		host:     host,
		username: username,
		password: password,
	}
}

type client struct {
	isAuthenticated bool
	token           string
	host            string
	username        string
	password        string
}

func (c *client) Login() error {
	basicAuth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.username, c.password)))
	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", c.host, loginUrl), nil)
	if err != nil {
		return fmt.Errorf("can't create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/jwt")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", basicAuth))
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("can't login: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Non-OK HTTP status: %d", response.StatusCode))
	}
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("can't read token: %w", err)
	}
	c.token = string(bodyBytes)
	c.isAuthenticated = true
	return nil
}

func (c *client) UploadArtifact(fileData []byte) (string, error) {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("artifact", string(fileData))
	writer.Close()

	headers := map[string][]string{
		"Accept":        {"application/json"},
		"Content-Type":  {writer.FormDataContentType()},
		"Authorization": {fmt.Sprintf("Bearer %s", c.token)},
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", c.host, uploadUrl), body)
	if err != nil {
		return "", fmt.Errorf("can't create request: %w", err)
	}
	req.Header = headers

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("can't send upload request: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusCreated {
		return "", errors.New(fmt.Sprintf("Non-OK HTTP status: %d", response.StatusCode))
	}
	location := response.Header.Get("Location")
	if location == "" {
		return "", errors.New("no location")
	}
	index := strings.LastIndex(location, "/")
	if index < 0 {
		return location, nil
	}
	id := location[index+1:]
	return id, nil
}
