package utils

import (
	"io"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

func BuildRequest(method string, url string, body io.Reader, headers http.Header) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if headers != nil {
		req.Header = headers
	}
	return req, nil
}

func CURLToByte(client http.Client, req *http.Request) ([]byte, error) {
	res, err := client.Do(req)
	defer func() {
		if res != nil && res.Body != nil {
			_ = res.Body.Close()
		}
	}()

	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func CURL(client http.Client, req *http.Request, out interface{}) error {
	b, err := CURLToByte(client, req)
	if err != nil {
		return err
	}
	if err = jsoniter.Unmarshal(b, &out); err != nil {
		return err
	}
	return err
}
