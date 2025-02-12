package sms

import (
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"ussd_ethereum/internal/service"
)

type SMSService struct {
	apiService *service.APIService
}

func NewSMSService(apiService *service.APIService) *SMSService {
	return &SMSService{
		apiService: apiService,
	}
}

func (s *SMSService) Send(message string, recipients []string, enqueue bool, sender_id string) (*http.Response, error) {
	URL := s.apiService.MakeURL("/version1/messaging")

	form := url.Values{
		"username":    {s.apiService.Username},
		"to":          recipients,
		"message":     {message},
		"bulkSMSMode": {"1"},
	}

	if sender_id != "" {
		form.Add("from", sender_id)
	}

	if enqueue {
		form.Add("enqueue", "1")
	}

	slog.Info("sms_send", "url", URL, "payload", form.Encode(), "headers", s.apiService.Headers)
	return s.apiService.MakeRequest(URL, "POST", s.apiService.Headers, strings.NewReader(form.Encode()), nil)

}
