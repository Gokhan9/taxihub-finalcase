package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Driver Client, API Gateway'in DriverService'e istek atmasını sağlayan HTTP Client
type DriverClient struct {
	baseURL string
	client  *http.Client
}

/*
en tamal amacı driver-service ile iletişim kurmak.(Get,post,put ve delete)
*/
func NewDriverClient() *DriverClient {

	// * Env'den DRIVER_SERVICE_URL'ini alır.
	baseURL := os.Getenv("DRIVER_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8081" //Localde geliştirme yapanın işini kolaylaştırmak için baseurl'i bastık.
	}

	return &DriverClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// TODO: HELPER: DriverService için bir "GET" isteği.
func (dc *DriverClient) Get(path string) (*http.Response, error) {

	url := fmt.Sprintf("%s%s", dc.baseURL, path)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return dc.client.Do(req)
}

// TODO: HELPER: "POST"
func (dc *DriverClient) Post(path string, body interface{}) (*http.Response, error) {

	jsonBody, _ := json.Marshal(body)

	url := fmt.Sprintf("%s%s", dc.baseURL, path)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return dc.client.Do(req)
}

// TODO: HELPER: PUT
func (dc *DriverClient) Put(path string, body interface{}) (*http.Response, error) {

	jsonBody, _ := json.Marshal(body)

	url := fmt.Sprintf("%s%s", dc.baseURL, path)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return dc.client.Do(req)
}

// TODO: HELPER: DELETE
func (dc *DriverClient) Delete(path string) (*http.Response, error) {

	url := fmt.Sprintf("%s%s", dc.baseURL, path)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	return dc.client.Do(req)
}

// TODO: HTTP'den dönen response'un body'sini okuyup str formatında geri döndürür.
func ReadResponseBody(resp *http.Response) (string, error) {

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
