package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/prometheus/common/log"
)

// response contains details of a response
type response struct {
	statusCode int
	size       int
}

// getResponse gets the response of a given url
func getResponse(url string) (*response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return &response{
		statusCode: resp.StatusCode,
		size:       len(body),
	}, nil
}

func main() {
	// list of urls from which response has to be obtained
	urls := []string{"https://www.google.com", "https://facebook.com", "https://twitt"}
	
	// final list of responses
	var responses []response

	var wg sync.WaitGroup
	wg.Add(len(urls))

	for _, url := range urls {
		go func(url string) {
			defer wg.Done()
			resp, err := getResponse(url)
			if err != nil {
				log.Error(err)
				return
			}
			responses = append(responses, *resp)
		}(url)
	}

	wg.Wait()

	// printing the final response
	fmt.Println(responses)
}
