package cloudflare

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// HTTPDo performs an HTTP request with the given method, URL, optional payload,
// and optional headers. Returns the response body as a trimmed string for any
// 2xx response, otherwise returns an error including a snippet of the body.
// If headers are provided, the first map will be used.
func HTTPDo(method, requestURL string, payload []byte, headers map[string]string) (string, error) {
	var bodyReader io.Reader
	if payload != nil {
		bodyReader = bytes.NewReader(payload)
	}

	request, err := http.NewRequest(method, requestURL, bodyReader)
	if err != nil {
		return "", err
	}

	if headers != nil {
		for k, v := range headers {
			request.Header.Set(k, v)
		}
	}
	if payload != nil && request.Header.Get("Content-Type") == "" {
		request.Header.Set("Content-Type", "application/json")
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode/100 != 2 {
		bodyBytes, _ := io.ReadAll(response.Body)
		return "", fmt.Errorf("unexpected status %d: %s", response.StatusCode, strings.TrimSpace(string(bodyBytes)))
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bodyBytes)), nil
}

// HTTPGet performs a simple HTTP GET request to the provided URL with optional headers.
func HTTPGet(requestURL string, headers map[string]string) (string, error) {
	return HTTPDo("GET", requestURL, nil, headers)
}

// HTTPPut performs a simple HTTP PUT to the provided URL with the given payload and optional headers.
func HTTPPut(requestURL string, payload []byte, headers map[string]string) (string, error) {
	return HTTPDo("PUT", requestURL, payload, headers)
}
