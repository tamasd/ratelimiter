// A simple rate limiter middleware.
// Copyright (c) 2020. Tamás Demeter-Haludka
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	url               = flag.String("url", "http://localhost:8080", "target url")
	count             = flag.Uint("count", 4096, "number of messages")
	concurrency       = flag.Uint("concurrency", 4, "concurrency")
	requestsPerSecond = flag.Uint("reqs", 100, "requests per second")
)

// worker is a worker thread that sends the requests in a given interval.
func worker(inputch chan func(), workers *sync.WaitGroup, requestsPerSecond uint) {
	defer workers.Done()
	reqtime := time.Second / time.Duration(requestsPerSecond)
	for job := range inputch {
		start := time.Now()
		job()
		time.Sleep(reqtime - time.Since(start))
	}
}

type result struct {
	statusCode int
	duration   time.Duration
}

// loadTest executes the load test.
func loadTest(url string, count, concurrency, requestsPerSecond uint) {
	var wg sync.WaitGroup

	inputch := make(chan func())

	wg.Add(int(concurrency))
	for i := uint(0); i < concurrency; i++ {
		go worker(inputch, &wg, requestsPerSecond)
	}

	client := &http.Client{}
	var resultmtx sync.Mutex
	results := make([]result, 0, int(count))

	start := time.Now()

	for i := uint(0); i < count; i++ {
		inputch <- func() {
			req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(""))
			start := time.Now()
			resp, err := client.Do(req)
			if err != nil {
				log.Println(err)
				return
			}
			duration := time.Since(start)

			// record the result
			resultmtx.Lock()
			results = append(results, result{
				statusCode: resp.StatusCode,
				duration:   duration,
			})
			resultmtx.Unlock()
		}
	}

	close(inputch)

	wg.Wait()

	duration := time.Since(start)

	summarize(results, duration)
}

// summarize prints a summary of the measurements.
func summarize(results []result, duration time.Duration) {
	var (
		successes, fails                           int
		allDuration, successDuration, failDuration time.Duration
	)

	for _, res := range results {
		if res.statusCode == http.StatusOK {
			successes++
			successDuration += res.duration
		} else {
			fails++
			failDuration += res.duration
		}
		allDuration += res.duration
	}

	seconds := float64(duration) / float64(time.Second)
	resultCount := len(results)

	fmt.Printf(
		"Average response time: %s, success: %s, failure: %s\n",
		allDuration/time.Duration(resultCount),
		successDuration/time.Duration(resultCount),
		failDuration/time.Duration(resultCount),
	)
	fmt.Printf("Duration: %s\n", duration)
	fmt.Printf("Requests per second: %.2f, success: %.2f\n",
		float64(resultCount)/seconds,
		float64(successes)/seconds,
	)
	fmt.Printf("Successes: %d\n", successes)
	fmt.Printf("Failures: %d\n", fails)
}

func main() {
	flag.Parse()
	loadTest(*url, *count, *concurrency, *requestsPerSecond)
}
