Starting point on my machine.

```sh
$ go test -bench=.

goos: darwin
goarch: arm64
pkg: metrics
BenchmarkMetrics/Average_age-8              1047           1141434 ns/op
BenchmarkMetrics/Average_payment-8            64          17101507 ns/op
BenchmarkMetrics/Payment_stddev-8             33          34107580 ns/op
PASS
ok      metrics 5.030s
```

After pulling out the ages into their own array using the "struct of arrays"
pattern:

```sh
BenchmarkMetrics/Average_age-8              2370            502892 ns/op
```
This is a 2.27x speed-up.

I'm not sure this optimization is allowed because it does open us up to
potential overflow, but doing only one division at the end gives us a
significant speed-up. Since we know the values are ages and therefore bounded
between [0, 120], we'd have to have an array of users orders of magnitude larger
than earth's population to overflow an `int64`, so I believe this is a safe
optimization.

```sh
BenchmarkMetrics/Average_age-8             12699             94534 ns/op
```

This is a 12.07x speed-up

We can go even one step further here and realize that the slice contains the
length of the array, so we can save an increment instruction in each loop
iteration.

```sh
BenchmarkMetrics/Average_age-8             31690             37873 ns/op
```

This is a 30.14x speed-up
