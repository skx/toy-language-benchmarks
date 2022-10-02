# Toy Language Benchmark

I dabble with toy scripting languages, having implemented several:

* TCL
* FORTH
* Lisp
* BASIC
* Monkey

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
BenchmarkGoFactorial-4    	 4519579	       260.1 ns/op
BenchmarkYALFactorial-4   	    1064	   1061797 ns/op
BenchmarkTCLFactorial-4   	      32	  38524019 ns/op
PASS
ok  	github.com/skx/toy-language-benchmarks	3.976s
```

Which means in terms of speed:

* Native Go
* Lisp
* TCL
