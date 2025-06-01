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
		fmt.Printf("Target: %s, Depth: %d, Status: %d, Latency: %v\n", target.Name, depth, result.Status, result.Latency)

		if result.Status != 200 || result.Latency > timeout {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func runSequentialTargets(targets []Target, maxDepth int, timeout time.Duration, cooldown time.Duration) {
	for i, target := range targets {
		fmt.Printf("\nTesting %s...\n", target.Name)

		// Create buffered channel with enough capacity for all possible results
		results := make(chan Result, maxDepth-8) // from 9 to maxDepth
		var wg sync.WaitGroup

		wg.Add(1)
		go testTarget(target, maxDepth, timeout, results, &wg)

		// Start a goroutine to wait for the test to complete
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
			close(results)
		}()

		// Wait for completion
		<-done

		if i < len(targets)-1 {
			fmt.Printf("\n⏳ Cooling down for %v before switching target...\n", cooldown)
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
