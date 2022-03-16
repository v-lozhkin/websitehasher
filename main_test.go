package main

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestScheduler_Run(t *testing.T) {
	cases := []struct {
		name      string
		responses []string
	}{
		{
			name:      "simple case with one",
			responses: []string{"Test response"},
		},
		{
			name:      "case with many",
			responses: []string{"Response one", "Кириллический ответ", "12311612", "<HTML>smth</HTML>"},
		},
		{
			name:      "empty response",
			responses: []string{""},
		},
	}

	for _, cse := range cases {
		cse := cse

		t.Run(cse.name, func(t *testing.T) {
			stubs := make([]*httptest.Server, 0, len(cse.responses))
			expected := make(map[string]string, len(cse.responses))

			for _, response := range cse.responses {
				response := response

				stub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					if _, err := fmt.Fprint(w, response); err != nil {
						t.Fatal(err)
					}
				}))
				stubs = append(stubs, stub)

				expected[stub.URL] = fmt.Sprintf("%x", md5.Sum([]byte(response)))
			}

			paths := make([]string, 0, len(cse.responses))
			for _, stub := range stubs {
				paths = append(paths, stub.URL)
			}

			parser := New(paths, 2)
			for _, result := range parser.Run() {
				if expected[result.path] != result.result {
					t.Fatalf(
						"wrong answer for %s, expected: %s, got: %s",
						result.path,
						expected[result.path],
						result.result,
					)
				}
			}
		})
	}
}
