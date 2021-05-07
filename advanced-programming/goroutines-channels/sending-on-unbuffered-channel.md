# Sending on an Unbuffered Channel That Blocks

Given this program:

```go
package main

import "time"

func main() {
  ch := make(chan int)

  go func() {
    // We will block here because there is no receive until the sleep is over
    ch <- 1
  }()

  time.Sleep(time.Duration(3) * time.Second)
  // Finally unblock the goroutine and receive
  <-ch
}
```

We want to understand what happens when the goroutine is blocked on a send.
