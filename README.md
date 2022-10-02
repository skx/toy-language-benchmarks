# Toy Language Benchmark

I dabble with toy scripting languages, having implemented several:

* [TCL](https://github.com/skx/critical)
* FORTH
* [Lisp](https://github.com/skx/yal)
* BASIC
* Monkey
* [evalfilter](https://github.com/skx/evalfilter)

And some other misc things.

This repository is designed to have a quick benchmark, unscientific, of those implementations.

## Benchmark

The idea is to calculate 100! (i.e. factorial of 100), and see how they compare.

To run the benchmark clone this repository then run:

    go test -run=Bench -bench=.

Or

    go test -run=Bench -bench=. -benchtime=30s

## Results

As of today the results, on my desktop system, look like this:

```
frodo ~/Repos/github.com/skx/toy-language-benchmarks $ go test -run=Bench -bench=. -benchtime=30s
goos: linux
goarch: amd64
pkg: github.com/skx/toy-language-benchmarks
cpu: AMD A10-6800K APU with Radeon(tm) HD Graphics
BenchmarkGoFactorial-4             141407989           249 ns/op
BenchmarkYALFactorial-4                32666       1175922 ns/op
BenchmarkTCLFactorial-4                  834      46042433 ns/op
BenchmarkEvalFilterFactorial-4        160035        227587 ns/op
PASS
ok  	github.com/skx/toy-language-benchmarks	192.035s
```

Which means in terms of speed:

* Native Go
* Evalfilter
  * Expected as this uses a bytecode virtual machine, which is pretty fast.
* Lisp
  * Expected as TCO makes this a reasonably fast benchmark.
* TCL
