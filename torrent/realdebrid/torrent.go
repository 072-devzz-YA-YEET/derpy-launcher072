package realdebrid

import (
	"fmt"
	"net/http"
)

func (client *RealDebridClient) AddTorrent() {

}

type AvailableHost struct {
	Host        string `json:"host"`
	MaxFileSize int    `json:"max_file_size"`
}

type AvailableHostsResponse struct {
	availableHosts []AvailableHost
}

func (client *RealDebridClient) AvailableHosts() ([]AvailableHost, error) {

	var availableHostResponse AvailableHostsResponse

	req, err := client.newRequest(http.MethodGet, "/torrents/availableHosts", nil, "", nil)
	if err != nil {
		return nil, fmt.Errorf("get request failed while requesting available hosts: %w", err)
	}
	var result []AvailableHost

	err = client.do(req, &result)
	if err != nil {
		return nil, err
	}

	for _, host := range availableHostResponse.availableHosts {
		result = append(result, host)
	}

	return result, nil
}
