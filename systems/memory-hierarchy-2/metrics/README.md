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
