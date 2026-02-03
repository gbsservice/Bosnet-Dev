package api

import (
	"crypto/tls"
	"time"

	"github.com/go-resty/resty/v2"
)

func Client() *resty.Client {
	return resty.New().SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).SetTimeout(120 * time.Second)
}
