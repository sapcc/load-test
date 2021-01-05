package loadtest

import (
	"net/http"

	"net/url"

	vegeta "github.com/tsenart/vegeta/lib"
)

func NewStaticTargeter(baseURL, endpoint, authToken string) (vegeta.Targeter, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	u, err = u.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	header := http.Header{}
	header.Add("X-Auth-Token", authToken)
	header.Add("Content-Type", "application/json")
	return vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    u.String(),
		Header: header,
	}), nil
}
