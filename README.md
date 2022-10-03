# Toy Language Benchmark

I dabble with toy scripting languages, having implemented several:

* [TCL](https://github.com/skx/critical)
* [FORTH](https://github.com/skx/foth)
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
$ go test -run=Bench -bench=.
goos: linux
goarch: amd64
pkg: github.com/skx/toy-language-benchmarks
cpu: AMD A10-6800K APU with Radeon(tm) HD Graphics
BenchmarkGoFactorial-4              4783036         252 ns/op
BenchmarkYALFactorial-4                1056     1102693 ns/op
BenchmarkTCLFactorial-4                  33    39381176 ns/op
BenchmarkEvalFilterFactorial-4         5461      204603 ns/op
BenchmarkFothFactorial-4              10000      145922 ns/op
PASS
ok    github.com/skx/toy-language-benchmarks  6.709s
```

Which means in terms of speed:

* Native Go
* FOTH
  * Expected as this is essentially a bytecode virtual machine, which is pretty fast.
* Evalfilter
  * Expected as this uses a bytecode virtual machine, which is pretty fast.
* Lisp
  * Expected as TCO makes this a reasonably fast benchmark.
* TCL
  * The constant string-conversion kills performance of maths.
