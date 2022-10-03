# Toy Language Benchmark

I dabble with toy scripting languages, having implemented several including:

* [TCL](https://github.com/skx/critical)
* [FORTH](https://github.com/skx/foth)
* [Lisp](https://github.com/skx/yal)
* [Monkey](https://github.com/skx/monkey)
* [evalfilter](https://github.com/skx/evalfilter)

This repository is designed to host a quick, unscientific, benchmark of those implementations.




## Benchmark

The idea is to calculate 100! (i.e. factorial of 100), and see how they compare.

To run the benchmark clone this repository then run:

    go test -run=Bench -bench=.

Or

    go test -run=Bench -bench=. -benchtime=30s

* **NOTE**: Some of the implementations actually fail!
  * (Specifically because they cannot handle numbers as large as those involved in 100!.)
  * I'm glossing over that because the benchmark does measure the overhead of recursing, numeric operations, and similar.




## Results

As of today the results, on my desktop system, look like this:

```sh
$ go test -run=Bench -bench=.
..
BenchmarkGoFactorial-10              1939436        620 ns/op
BenchmarkEvalFilterFactorial-10        13442      89138 ns/op
BenchmarkFothFactorial-10              10000     120620 ns/op
BenchmarkMonkeyFactorial-10             5034     231593 ns/op
BenchmarkYALFactorial-10                3440     327744 ns/op
BenchmarkTCLFactorial-10                 100   11065288 ns/op
PASS
ok   github.com/skx/toy-language-benchmarks 8.841s
```

Which means in terms of speed:

* Native Go
  * Obviously this would be best, being compiled.
* Evalfilter
  * Expected as this uses a bytecode virtual machine, which is pretty fast.
* FOTH
  * Expected as this is essentially a bytecode virtual machine, which is pretty fast.
* Monkey
  * This was a surprise, as there's no particular optimisation or cleverness here.
* Lisp
  * Expected as TCO makes this a reasonably fast benchmark.
* TCL
  * The constant string-conversion kills performance of maths.
