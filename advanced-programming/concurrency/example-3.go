package main

import (
	"fmt"
	"net/http"
	"sync"
)

func main_orig() {
	var urls = []string{
		"https://bradfieldcs.com/courses/architecture/",
		"https://bradfieldcs.com/courses/networking/",
		"https://bradfieldcs.com/courses/databases/",
	}
	var wg sync.WaitGroup
	for i := range urls {
		go func(i int) {
			wg.Add(1)
			// Decrement the counter when the goroutine completes.
			defer wg.Done()

			_, err := http.Get(urls[i])
			if err != nil {
				panic(err)
			}

			fmt.Println("Successfully fetched", urls[i])
		}(i)
	}

	// Wait for all url fetches
	wg.Wait()
	fmt.Println("all url fetches done!")
}

// The problem here is that we're putting the wg.Add(1) inside the goroutine so
// we can't be sure that the wait happens after the add
func main() {
	var urls = []string{
		"https://bradfieldcs.com/courses/architecture/",
		"https://bradfieldcs.com/courses/networking/",
		"https://bradfieldcs.com/courses/databases/",
	}
	var wg sync.WaitGroup
	for i := range urls {
		wg.Add(1)
		go func(i int) {
			// Decrement the counter when the goroutine completes.
			defer wg.Done()

			_, err := http.Get(urls[i])
			if err != nil {
				panic(err)
			}

			fmt.Println("Successfully fetched", urls[i])
		}(i)
	}

	// Wait for all url fetches
	wg.Wait()
	fmt.Println("all url fetches done!")
}
