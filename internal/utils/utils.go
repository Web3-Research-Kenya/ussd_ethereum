package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"strings"
)

func GetLastIdx(input string) string {
	return "a"
}

func FormatPhoneNumber(phone string) (string, error) {
	if strings.HasPrefix(phone, "+254") && len(phone) == 13 {
		return phone, nil
	}

	if len(phone) == 10 && strings.HasPrefix(phone, "0") {
		internationalFormat := "+254" + phone[1:]
		return internationalFormat, nil
	}

	return "", fmt.Errorf("invalid phone number format")
}

func WeiToEth(wei *big.Int) *big.Float {
	// 1 ETH = 10^18 Wei
	eth := new(big.Float).SetInt(wei)
	return new(big.Float).Quo(eth, big.NewFloat(1e18))
}

func MapToReader(m map[string]any) (io.Reader, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(b))

	return bytes.NewReader(b), nil
}
