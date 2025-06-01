package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Target struct {
	BaseURL string
	Name    string
}

type Result struct {
	Depth   int
	Latency time.Duration
	Status  int
	Error   string
}

func sendRequest(target Target, depth int) Result {
	start := time.Now()
	resp, err := http.Get(target.BaseURL + fmt.Sprintf("%d", depth))
	latency := time.Since(start)

	if err != nil {
		return Result{
			Depth:   depth,
			Latency: latency,
			Status:  0,
			Error:   err.Error(),
		}
	}
	defer resp.Body.Close()

	return Result{
		Depth:   depth,
		Latency: latency,
		Status:  resp.StatusCode,
		Error:   "",
	}
}

func testTarget(target Target, maxDepth int, timeout time.Duration, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for depth := 9; depth <= maxDepth; depth++ {
		result := sendRequest(target, depth)
		results <- result

		if result.Status != 200 || result.Latency > timeout {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	close(results)
}

func runSequentialTargets(targets []Target, maxDepth int, timeout time.Duration, cooldown time.Duration) {
	for i, target := range targets {
		results := make(chan Result, 10)
		var wg sync.WaitGroup

		wg.Add(1)
		go testTarget(target, maxDepth, timeout, results, &wg)

		wg.Wait()

		if i < len(targets)-1 {
			fmt.Println("⏳ Cooling down before switching target...")
			time.Sleep(cooldown)
		}
	}
}

func main() {
	targets := []Target{
		{BaseURL: "http://localhost:3000/api/image/", Name: "Bun Service"},
		{BaseURL: "http://localhost:3001/api/image/", Name: "Node Service"},
		{BaseURL: "http://localhost:3002/api/image/", Name: "Deno Service"},
	}

	maxDepth := 20
	timeout := 1000000000 * time.Millisecond
	cooldown := 5 * time.Second

	runSequentialTargets(targets, maxDepth, timeout, cooldown)
}
