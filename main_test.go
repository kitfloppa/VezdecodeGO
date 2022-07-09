package main

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"testing"
)

func TestServerOneThreadAsync(t *testing.T) {
	go New()

	data := url.Values{
		"timeDuration": {"0h0m10s"},
		"async":        {""},
	}

	for i := 0; i < 5; i++ {
		_, err := http.PostForm("http://localhost:8081/add", data)
		if err != nil {
			log.Fatal(err)
		}
	}

	resp, err := http.Get("http://localhost:8081/time")

	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	require.Equal(t, "40s", string(body))
}

func TestServerOneThreadSync(t *testing.T) {
	go New()

	data := url.Values{
		"timeDuration": {"0h0m10s"},
		"sync":         {""},
	}

	for i := 0; i < 2; i++ {
		_, err := http.PostForm("http://localhost:8081/add", data)
		if err != nil {
			log.Fatal(err)
		}
	}

	resp, err := http.Get("http://localhost:8081/time")

	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	require.Equal(t, "0s", string(body))
}
