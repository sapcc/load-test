package loadtest

import (
	"net/http"

	vegeta "github.com/tsenart/vegeta/lib"
)

func NewStaticTargeter(ctx *Context) vegeta.Targeter {
	header := http.Header{}
	header.Add("X-Auth-Token", ctx.AuthToken)
	header.Add("Content-Type", "application/json")

	return vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    ctx.Url,
		Header: header,
	})
}
