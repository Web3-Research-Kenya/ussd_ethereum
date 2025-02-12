package service

import (
	"io"
	"log/slog"
	"net/http"
	"strings"
)

const (
	PRODUCTION_DOMAIN = "africastalking.com"
	SANDBOX_DOMAIN    = "sandbox.africastalking.com"
)

var (
	BASE_URL = "https://api." + PRODUCTION_DOMAIN
)

type APIService struct {
	Username string
	APIKey   string
	BaseUrl  string
	Headers  http.Header
}

func NewAPIService(username, apiKey string) *APIService {
	var baseUrl string = "https://api." + PRODUCTION_DOMAIN
	if username == "sandbox" {
		baseUrl = "https://api." + SANDBOX_DOMAIN
	}
	// headers := &http.Header{
	// 	"Accept":       []string{"application/json"},
	// 	"Content-Type": []string{"application/x-www-form-urlencoded"},
	// 	"apiKey":       []string{apiKey},
	// }
	headers := http.Header{}
	headers.Add("Accept", "application/json")
	headers.Add("Content-Type", "application/x-www-form-urlencoded")
	headers.Add("apiKey", apiKey)

	return &APIService{
		Username: username,
		APIKey:   apiKey,
		BaseUrl:  baseUrl,
		Headers:  headers,
	}
}

func (s *APIService) IsSandbox() bool {
	return s.Username == "sandbox"
}

func (s *APIService) MakeURL(path string) string {
	if s.IsSandbox() {
		return "https://api." + SANDBOX_DOMAIN + path
	} else {
		return "https://api." + PRODUCTION_DOMAIN + path
	}
}

func (s *APIService) MakeGetRequest(url string, headers http.Header, data io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest("GET", url, data)
	if err != nil {
		return nil, err
	}
	s.addRequestHeaders(req, headers)
	return req, nil
}

func (s *APIService) MakePostRequest(url string, headers http.Header, data io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest("POST", url, data)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("apiKey", s.APIKey)

	return req, nil
}

func (s *APIService) MakeRequest(url, method string, headers http.Header, data io.Reader, params []string) (response *http.Response, err error) {
	method = strings.ToUpper(method)
	client := http.Client{}

	switch method {
	case "GET":
		request, err := s.MakeGetRequest(url, headers, data)
		if err != nil {
			return nil, err
		}
		response, err = client.Do(request)
		if err != nil {
			return nil, err
		}
	case "POST":
		request, err := s.MakePostRequest(url, headers, data)
		if err != nil {
			return nil, err
		}
		response, err = client.Do(request)
		if err != nil {
			return nil, err
		}
	}
	return response, nil
}

func (s *APIService) addRequestHeaders(req *http.Request, headers http.Header) {
	for key, values := range s.Headers {
		for _, value := range values {
			req.Header.Set(key, value)
		}
		slog.Info("add_request_headers", "api_headers", key)
	}

	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

}
