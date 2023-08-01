package services

import (
	"github.com/robfig/cron/v3"
	"log"
	"sano/config"
	"sano/database"
	"strings"
)

func isAvailableProtocol(protocol string) bool {
	availableProtocols := []string{protocolHttp, protocolHttps}
	available := false

	for _, availableProtocol := range availableProtocols {
		if availableProtocol == protocol {
			available = true
		}
	}

	return available
}

func RunLookup(services []config.Service, defaultCronTiming string, c *cron.Cron) {
	for _, service := range services {
		cronTime := defaultCronTiming
		if service.Cron != nil {
			cronTime = *service.Cron
		}

		protocol := service.Url[:strings.Index(service.Url, "://")]
		if !isAvailableProtocol(protocol) {
			log.Printf("I don't know which protocol you want to use for [%s], please fix it or maybe it's unknown for me then open an issue pls.", service.Name)
			continue
		}

		cronFunc := func() {
			var err error

			switch protocol {
			case protocolHttp:
				err = runHealthCheckHttp(service)
			case protocolHttps:
				err = runHealthCheckHttp(service)
			default:
				log.Printf("I don't know which protocol you want to use for [%s], please fix it or maybe it's unknown for me then open an issue pls.", service.Name)
			}

			status := err == nil
			if !status {
				log.Printf("[%s] is offline, message: %s", service.Name, err)
			}

			database.StoreLookup(service, status)
		}
		_, err := c.AddFunc(cronTime, cronFunc)
		if err != nil {
			log.Panicln(err)
		}
	}
}
