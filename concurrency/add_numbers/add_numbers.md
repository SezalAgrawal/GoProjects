# Observations

## Concurrency without parallelism

### Command

```sh
GOGC=off go test -cpu 1 -run none -bench . -benchtime 3s
```

### Response

```sh
Processing 10000000 numbers using 4 goroutines
goos: darwin
goarch: amd64
pkg: github.com/goProjects/concurrency/add_numbers
BenchmarkSequential                  564           7017853 ns/op
BenchmarkConcurrent                  444           9672264 ns/op
BenchmarkSequentialAgain             552           6398380 ns/op
BenchmarkConcurrentAgain             492           7722396 ns/op
PASS
ok      github.com/goProjects/concurrency/add_numbers   25.367s
```

### Conclusion

Concurrent takes more time as add_numbers is a cpu bound process
and lot of context switch overhead takes place as there is only single thread.

------------

## Concurrency with parallelism

### Command

```sh
GOGC=off go test -run none -bench . -benchtime 3s
```

### Response

```sh
Processing 10000000 numbers using 4 goroutines
goos: darwin
goarch: amd64
pkg: github.com/goProjects/concurrency/add_numbers
BenchmarkSequential-4                561           6213913 ns/op
BenchmarkConcurrent-4               1034           3231111 ns/op
BenchmarkSequentialAgain-4           562           6507637 ns/op
BenchmarkConcurrentAgain-4          1107           3410508 ns/op
PASS
ok      github.com/goProjects/concurrency/add_numbers   17.052s
```

### Conclusion

Concurrent takes very less time than sequential as 4 goroutines are running
in parallel and thus completing their work concurrently.
