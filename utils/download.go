package utils

import (
	"log"
	"net/http"
	"net/url"
	"os"
)

func DownloadFile(uri string, fileName string) {
	_, err := url.Parse(uri)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	resp, err := client.Get(uri)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}
