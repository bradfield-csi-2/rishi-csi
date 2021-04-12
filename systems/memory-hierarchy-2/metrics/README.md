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

Now onto the payment average and standard deviation.

The first optimization that comes to mind is to treat the payment amount as
cents only instead of dollars and cents. This also as the nice outcome of
avoiding floating point rounding errors in general. Store money in ints until
the very last moment.

```sh
BenchmarkMetrics/Average_payment-8           172           6828006 ns/op
BenchmarkMetrics/Payment_stddev-8             81          14325964 ns/op
```

This is a 2.50x speed-up in average and a 2.38x speed-up in standard deviation.

Just as we did with users, we can add the payment information to our struct of
arrays. Since we have no need to group payments by user for this calculation, we
can do away with the nested loop and keep all the amounts in just one array.

```sh
BenchmarkMetrics/Average_payment-8           235           5091768 ns/op
BenchmarkMetrics/Payment_stddev-8            184           6461564 ns/op
```
This is a 3.36x speed-up in average and a 5.28x speed-up in standard deviation.

A tiny improvement to standard deviation might be to skip the division by 100 in
each iteration and save it until the end, but this had a negligible impact.

```sh
BenchmarkMetrics/Payment_stddev-8            186           6427847 ns/op
```

Another tiny improvement is not incrementing the count in standard deviation
```sh
BenchmarkMetrics/Payment_stddev-8            181           6402519 ns/op
```

A potentially more significant improvement is to use 32-bit integers in our
payments array instead of 64-bit as we know the bounds of each payment falls
well below 2^32 - 1 cents. This will double the cache hit rate, but there was no
significant speed-up.
```sh
BenchmarkMetrics/Average_payment-8           237           5024183 ns/op
BenchmarkMetrics/Payment_stddev-8            189           6294602 ns/op
```

Again as with users, assuming we can safely avoid overflow, summing the total
payments instead of computing the running in average in each iteration will
provide a big boost to speed. Since we have 10^6 payments of at most 10^8 cents,
the sum of payments will be at most 10^14 cents, which is within 2^64 - 1 ~
10^19 of a uint64 accumulator.
```sh
BenchmarkMetrics/Average_payment-8          3415            342089 ns/op
BenchmarkMetrics/Payment_stddev-8            756           1598360 ns/op
```
This is a 49.99x speed-up in average and 21.34x speed-up in standard deviation.

A nice further improvement gain is to unroll the loop. Unrolling in the average
payment calculation gives
```sh
BenchmarkMetrics/Average_payment-8          4118            292331 ns/op
BenchmarkMetrics/Payment_stddev-8            765           1565053 ns/op
```
This is a 58.50x speed-up in average and 21.79x speed-up in standard deviation.
