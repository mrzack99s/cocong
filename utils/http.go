package utils

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/mrzack99s/cocong/types"
)

func HttpPing(method, fullUrl string) bool {
	req, err := http.NewRequest(method, fullUrl, nil)
	if err != nil {
		return false
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	resp.Body.Close()

	return resp.StatusCode == 200
}

func HttpPingWithCheckStatusCode(method, fullUrl string, checkWithStatusCode int) bool {

	req, err := http.NewRequest(method, fullUrl, nil)
	if err != nil {
		return false
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	resp.Body.Close()

	if checkWithStatusCode == 0 && resp.StatusCode == 200 {
		return true
	}

	return checkWithStatusCode == resp.StatusCode
}

func HttpJSONRequestWithBytesResponse(method, fullURL, token string, requestData io.Reader) ([]byte, error) {

	req, err := http.NewRequest(method, fullURL, requestData)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("api-token", token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return respbody, errors.New("response code is not ok")
	}

	return respbody, nil
}

func HttpRequestWithBytesResponse(r types.HttpRequestType) ([]byte, error) {

	reader := bytes.NewReader(r.Data)

	req, err := http.NewRequest(r.Method, r.FullURL, reader)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	if r.ContentType != "" {
		req.Header.Set("Content-Type", r.ContentType)
	}
	for _, header := range r.HeaderAdditional {
		req.Header.Set(header.Name, header.Value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	respbody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, errors.New(string(respbody))
	}

	return respbody, nil
}
