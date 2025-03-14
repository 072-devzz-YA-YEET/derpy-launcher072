package realdebrid

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

var (
	ErrorInvalidRequest  = errors.New("invalid request")
	ErrorInvalidURL      = errors.New("invalid URL")
	ErrorCannotParsePath = errors.New("cannot parse path")
	ErrorCannotReadBody  = errors.New("Cannot read body")
	Error401             = errors.New("Unauthorized")
	Error403             = errors.New("Forbidden")
	Error404             = errors.New("Not Found")
	Error500             = errors.New("Internal Server Error")
)

type RealDebridClient struct {
	ApiKey string
	client *http.Client
	BaseURL string
}

func NewRealDebridClient(apiKey string) *RealDebridClient {
	return &RealDebridClient{
		ApiKey: apiKey,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		BaseURL: "https://api.real-debrid.com/rest/1.0",
	}
}

func (client *RealDebridClient) newRequest(method, path string, headers http.Header, params string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, client.BaseURL+path+params, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+client.ApiKey)
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header[k] = v
	}
	return req, nil
}

func (client *RealDebridClient) do(req *http.Request, v interface{}) error {
	resp, err := client.client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotFound:
		return Error404
	case http.StatusUnauthorized:
		return Error401
	case http.StatusForbidden:
		return Error403
	case http.StatusInternalServerError:
		return Error500
	}

	defer resp.Body.Close()

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
		if err != nil {
			return err
		}
	}

	return nil
}