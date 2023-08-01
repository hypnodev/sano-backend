package services

import (
	"errors"
	"net/http"
	"sano/config"
)

const protocolHttp = "http"
const protocolHttps = "https"

func runHealthCheckHttp(service config.Service) error {
	response, err := http.Get(service.Url)
	if err != nil {
		return errors.New("connection to service failed")
	}

	if response.StatusCode >= 200 && response.StatusCode < 400 {
		return nil
	}

	return errors.New(response.Status)
}
