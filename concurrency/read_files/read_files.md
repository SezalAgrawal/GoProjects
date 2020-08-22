# Observations

## Concurrency without parallelism

### Command

```sh
GOGC=off go test -cpu 1 -run none -bench . -benchtime 3s
```

### Response

```sh
Processing 1000 files using 4 goroutines
goos: darwin
goarch: amd64
pkg: github.com/goProjects/concurrency/read_files
BenchmarkSequential                    3        1327811995 ns/op
BenchmarkConcurrent                   10         341825887 ns/op
BenchmarkSequentialAgain               3        1338574065 ns/op
BenchmarkConcurrentAgain              10         327539234 ns/op
PASS
ok      github.com/goProjects/concurrency/read_files    23.853s
```

### Conclusion

Concurrent process is almost 87% faster since context switch naturally happens.

------------

## Concurrency with parallelism

### Command

```sh
GOGC=off go test -run none -bench . -benchtime 3s
```

### Response

```sh
Processing 1000 files using 4 goroutines
goos: darwin
goarch: amd64
pkg: github.com/goProjects/concurrency/read_files
BenchmarkSequential-4                  3        1363172690 ns/op
BenchmarkConcurrent-4                 12         337046704 ns/op
BenchmarkSequentialAgain-4             3        1373322614 ns/op
BenchmarkConcurrentAgain-4             9         339966704 ns/op
PASS
ok      github.com/goProjects/concurrency/read_files    24.609s
```

### Conclusion

IO bound work works efficiently using goroutines on one thread, since context switch naturally happens. More hardware threads does not add any advantage.
