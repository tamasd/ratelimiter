[![Build Status](https://travis-ci.org/tamasd/ratelimiter.svg?branch=v1)](https://travis-ci.org/tamasd/ratelimiter)
[![codecov](https://codecov.io/gh/tamasd/ratelimiter/branch/v1/graph/badge.svg)](https://codecov.io/gh/tamasd/ratelimiter)
[![Go Report Card](https://goreportcard.com/badge/github.com/tamasd/ratelimiter)](https://goreportcard.com/report/github.com/tamasd/ratelimiter)
[![CodeFactor](https://www.codefactor.io/repository/github/tamasd/ratelimiter/badge)](https://www.codefactor.io/repository/github/tamasd/ratelimiter)
[![Known Vulnerabilities](https://snyk.io/test/github/tamasd/ratelimiter/v1/badge.svg)](https://snyk.io/test/github/tamasd/ratelimiter)

# Rate limiter middleware

## Manual testing

To start the server:

```
$ HOST=localhost go run example/server/server.go
```

The server uses the `HOST` and `PORT` environment variables. For more options,
run `go run example/server/server.go -help`.

To start the client:

```
$ go run example/client/client.go
```

For more options, run `go run example/client/client.go -help`.

## Benchmark result

This is the reason why the mutex-based bucket is the default:

```
ratelimiter/internal/bucket $ go test -bench . -cpu 1,2,4,8,16 -benchmem
goos: windows
goarch: amd64
pkg: github.com/tamasd/ratelimiter/internal/bucket
BenchmarkMutex          33322502                36.3 ns/op             0 B/op          0 allocs/op
BenchmarkMutex-2         8277094               140 ns/op               0 B/op          0 allocs/op
BenchmarkMutex-4         5529948               216 ns/op               0 B/op          0 allocs/op
BenchmarkMutex-8         5911280               205 ns/op               0 B/op          0 allocs/op
BenchmarkMutex-16        6629870               181 ns/op               0 B/op          0 allocs/op
BenchmarkChannel         1242219               948 ns/op              96 B/op          1 allocs/op
BenchmarkChannel-2        705910              1453 ns/op              96 B/op          1 allocs/op
BenchmarkChannel-4        999915              1118 ns/op              96 B/op          1 allocs/op
BenchmarkChannel-8       1000000              1057 ns/op              96 B/op          1 allocs/op
BenchmarkChannel-16      1000000              1037 ns/op              96 B/op          1 allocs/op
PASS
ok      github.com/tamasd/ratelimiter/internal/bucket   13.246s

```

The mutex-based approach is much faster and does not allocate memory.

Test machine: Windows 10, AMD 1800X, 32 GB 3600 MHz CL14.

## TODO

* consider moving the examples to a separate project to lower the number of
  dependencies
* implement load-based randomized shedding
