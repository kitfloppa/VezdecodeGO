package main

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"
)

func TestServerOneThreadAsync(t *testing.T) {
	httpServerExitDone := &sync.WaitGroup{}

	httpServerExitDone.Add(1)

	term := make(chan byte)

	go New(term, httpServerExitDone)
	time.Sleep(time.Millisecond) // sleep for server starting

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

	term <- 1

	httpServerExitDone.Wait()
}

func TestServerOneThreadSync(t *testing.T) {

	httpServerExitDone := &sync.WaitGroup{}

	httpServerExitDone.Add(1)

	term := make(chan byte)

	go New(term, httpServerExitDone)
	time.Sleep(time.Millisecond) // sleep for server starting

	data := url.Values{
		"timeDuration": {"0h0m10s"},
		"sync":         {""},
	}

	for i := 0; i < 1; i++ {
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

	term <- 1

	httpServerExitDone.Wait()
}

func TestServerStressAsync(t *testing.T) {

	httpServerExitDone := &sync.WaitGroup{}

	httpServerExitDone.Add(1)

	term := make(chan byte)

	go New(term, httpServerExitDone)
	time.Sleep(time.Millisecond) // sleep for server starting

	data := url.Values{
		"timeDuration": {"0h0m10s"},
		"async":        {""},
	}

	f := func(res chan<- byte) {
		for i := 0; i < 10; i++ {
			_, err := http.PostForm("http://localhost:8081/add", data)
			if err != nil {
				log.Fatal(err)
			}
		}
		res <- 1
	}

	results := make(chan byte, 100)
	for j := 0; j < 50; j++ {
		go f(results)
	}

	for j := 0; j < 50; j++ {
		<-results
	}

	resp, err := http.Get("http://localhost:8081/time")

	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	require.Equal(t, "1h23m10s", string(body))

	term <- 1

	httpServerExitDone.Wait()
}
