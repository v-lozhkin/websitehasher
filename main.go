package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
)

type Parser struct {
	workers chan struct{}
	paths   []string
}

type Parsed struct {
	path   string
	result string
}

func (s Parser) Run() []Parsed {
	wg := sync.WaitGroup{}
	result := make([]Parsed, len(s.paths))

	for i, path := range s.paths {
		// pin variables
		path := path
		i := i

		// wait for free worker
		<-s.workers

		wg.Add(1)
		go func() {
			defer func() {
				// release the worker
				s.workers <- struct{}{}
				wg.Done()
			}()

			result[i] = Parsed{
				path: path,
			}

			hash, err := getMD5Hash(path)
			if err != nil {
				result[i].result = fmt.Sprintf("failed to get result of response of %s: %s", path, err.Error())
			} else {
				result[i].result = hash
			}
		}()
	}
	// wait for all workers
	wg.Wait()

	return result
}

func getMD5Hash(path string) (string, error) {
	parsedURL, err := url.Parse(path)
	if err != nil {
		return "", fmt.Errorf("failed to parse path %s: %w", path, err)
	}

	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "https"
	}

	res, err := http.Get(parsedURL.String())
	if err != nil {
		return "", fmt.Errorf("failed to get resource %s: %w", parsedURL.String(), err)
	}

	defer func() {
		if err = res.Body.Close(); err != nil {
			fmt.Println("failed to close response body", err)
		}
	}()

	resp, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return fmt.Sprintf("%x", md5.Sum(resp)), nil
}

func New(paths []string, workers int) Parser {
	// init parser
	parser := Parser{
		workers: make(chan struct{}, workers),
		paths:   paths,
	}

	// load up workers
	for i := 0; i < workers; i++ {
		parser.workers <- struct{}{}
	}

	return parser
}

func main() {
	parallel := flag.Int("parallel", 10, "limit of a parallel requests")
	flag.Parse()

	paths := flag.Args()

	parser := New(paths, *parallel)
	for _, result := range parser.Run() {
		fmt.Printf("%s: %s\n", result.path, result.result)
	}
}
