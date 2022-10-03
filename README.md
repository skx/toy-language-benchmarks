# Toy Language Benchmark

I dabble with toy scripting languages, having implemented several including:

* [BASIC](https://github.com/skx/gobasic)
* [FORTH](https://github.com/skx/foth)
* [Lisp](https://github.com/skx/yal)
* [Monkey](https://github.com/skx/monkey)
* [TCL](https://github.com/skx/critical)
* [evalfilter](https://github.com/skx/evalfilter)

This repository is designed to host a quick, unscientific, benchmark of those implementations.




## Benchmark

The idea is to calculate 12! (i.e. factorial of 12), and see how they compare.

To run the benchmark clone this repository then run:

    go test -run=Bench -bench=.

Or

    go test -run=Bench -bench=. -benchtime=30s




## Results

As of today the results, on my desktop system, look like this:

```sh
$ go test -run=Bench -bench=.
goos: linux
goarch: amd64
pkg: github.com/skx/toy-language-benchmarks
cpu: AMD A10-6800K APU with Radeon(tm) HD Graphics
BenchmarkGoFactorial-4           35637698        30 ns/op
BenchmarkEvalFilterFactorial-4      61542     17458 ns/op
BenchmarkFothFactorial-4            44751     26275 ns/op
BenchmarkBASICFactorial-4           36735     32090 ns/op
BenchmarkMonkeyFactorial-4          14446     85061 ns/op
BenchmarkYALFactorial-4              2607    456757 ns/op
BenchmarkTCLFactorial-4               292   4085301 ns/op
PASS
ok  	github.com/skx/toy-language-benchmarks	21.479s
```

Which means in terms of speed:

* Native Go
  * Obviously this would be best, being compiled.
* Evalfilter
  * Expected as this uses a bytecode virtual machine, which is pretty fast.
* FOTH
  * Expected as this is essentially a bytecode virtual machine, which is pretty fast.
* BASIC
  * This is a real surprise, although the code is nice and simple so maybe it shouldn't be?
* Monkey
  * This was a surprise, as there's no particular optimisation or cleverness here.
* Lisp
  * Expected as TCO makes this a reasonably fast benchmark.
* TCL
  * The constant string-conversion kills performance of maths.
