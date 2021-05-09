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

Starting on the first line of the program, we create a channel `ch`. This is
actually a pointer to an `hchan` struct.

We the goroutine sends on `ch`, it calls then `chansend1` function in the
runtime with arguments of the channel itself (the `hchan` pointer) and an
`unsafe.Pointer` of the data (in our case the `int` 1.

This is a wrapper around `chansend`, which provides a couple extra arguments:
the value `true` for the `block` argument and the `callerpc`

After checking if the channel is closed or not ready, the `hchan` struct itself
is locked.

Then we check if there is a receive that we can message by dequeueing the `recvq` on the `hchan`. Since there is nothing on that queue yet in our program, we fail this conditional and move on.

The next check is if the `hchan`'s `qcount` is less than its `dataqsiz`. This
checking if there is space available to enqueue the send. In this case, since
our channel is unbuffered, we have 0 `dataqsiz` so we move on.

Since our `block` argument was `true`, that means we do want to block and so we
skip an early return and get to the blocking code.

We set up `gp` and `mysg`, which are pointers to a "g", or a goroutine.
