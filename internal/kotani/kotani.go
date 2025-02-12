package kotani

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	baseUrl = "https://sandbox-api.kotanipay.io/api/v3/customer/mobile-money"
)

type Client struct {
	apiKey string
}

func NewClient(apiKey string) *Client {

	return &Client{
		apiKey: apiKey,
	}
}

func (c *Client) CreateMobileCustomer(phone, country_code, network string) error {
	payload := strings.NewReader("{\"country_code\":\"EG\",\"phone_number\":\"\\\"\",\"network\":\"MPESA\",\"account_name\":\"test\"}")

	req, _ := http.NewRequest("POST", baseUrl, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))
	return nil
}
