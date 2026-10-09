package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"roof/vpos/repository"
	"time"
)

func PrepareRequest(b *repository.Bolt, path string, body []byte) (*http.Request, error) {
	baseURL := b.ConfigRepo.GetBaseURL()
	hreq, err := http.NewRequest("POST", baseURL+path, bytes.NewBuffer(body))

	if err != nil {
		return &http.Request{}, err
	}
	hreq.Header, err = CalculateSignature(string(body), b)
	if err != nil {
		return &http.Request{}, err
	}
	return hreq, nil

}

func SendRequest(req *http.Request) ([]byte, http.Header, error) {
	client := &http.Client{
		Timeout: time.Second * 30,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 400 {
		return nil, resp.Header, fmt.Errorf("request failed with status code %s", resp.Status)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return respBody, resp.Header, nil
}
