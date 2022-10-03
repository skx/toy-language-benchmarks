package main

import (
	"fmt"
	"testing"

	// critical
	"github.com/skx/critical/interpreter"
	criticalSTDLIB "github.com/skx/critical/stdlib"

	// evalfilter
	"github.com/skx/evalfilter/v2"

	// foth
	fothEval "github.com/skx/foth/foth/eval"

	// yal
	"github.com/skx/yal/builtins"
	"github.com/skx/yal/env"
	"github.com/skx/yal/eval"
	"github.com/skx/yal/primitive"
	"github.com/skx/yal/stdlib"
)

// fact is a benchmark implementation in pure-go for comparison purposes.
func fact(n int64) int64 {
	if n == 0 {
		return 1
	}
	return n * fact(n-1)
}

// BenchmarkGoFactorial allows running the golang benchmark.
func BenchmarkGoFactorial(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fact(100)
	}
}

// BenchmarkYALFactorial allows running the lisp benchmark.
func BenchmarkYALFactorial(b *testing.B) {

	// Create a new environment
	yalEnv := env.New()

	// Populate with the default primitives
	builtins.PopulateEnvironment(yalEnv)

	// The script we're going to run
	content := `
(define fact (lambda (n)
  (if (<= n 1)
    1
      (* n (fact (- n 1))))))

(fact 100)
`

	// Read the standard library
	pre := stdlib.Contents()

	// Prepend that to the users' script
	src := string(pre) + "\n" + string(content)

	// Create a new interpreter with that source
	yalInter := eval.New(src)

	// The return value
	var out primitive.Primitive

	// Run the benchmark
	for i := 0; i < b.N; i++ {

		// Run 100!
		out = yalInter.Evaluate(yalEnv)
	}

	// Did we get an error?  Then show it.
	if _, ok := out.(primitive.Error); ok {
		fmt.Printf("Error running: %v\n", out)
	}
}

// BenchmarkTCLFactorial allows running the TCL benchmark.
func BenchmarkTCLFactorial(b *testing.B) {

	input := `
proc fact {n} {
    if  {<= $n 1} {
        return 1
    } else {
        return [* $n [fact [- $n 1]]]
    }
}
fact 100
`
	// Load the standard library
	stdlib := criticalSTDLIB.Contents()

	// Create an interpreter with our source and stdlib.
	tclInter, err := interpreter.New(string(stdlib) + "\n" + input)
	if err != nil {
		fmt.Printf("initTCL; Error creating interpreter %s\n", err)
		panic(err)
		return
	}

	var out any

	// Run the script
	for i := 0; i < b.N; i++ {

		// Run 100!
		out, err = tclInter.Evaluate()

	}

	// Ensure the result was "ok".
	if err != nil && err != interpreter.ErrReturn {
		fmt.Printf("Error running program:%s\n", err)
		return
	}
	if out != "93326215443944102188325606108575267240944254854960571509166910400407995064242937148632694030450512898042989296944474898258737204311236641477561877016501813248.000000" {
		fmt.Printf("Unexpected result:%s\n", out)
	}
}

// BenchmarkEvalFilterFactorial allows running the evalfilter benchmark.
func BenchmarkEvalFilterFactorial(b *testing.B) {

	prg := `
function fact( n ) {
  if ( n <= 1 ) { return 1; }
  return ( n * fact( n - 1 ) );
}
fact(100);
return false;
`
	eval := evalfilter.New(string(prg))

	err := eval.Prepare([]byte{})
	if err != nil {
		fmt.Printf("Error compiling:%s\n", err.Error())
		panic(err)
	}

	for i := 0; i < b.N; i++ {

		_, err = eval.Execute(true)
		if err != nil {
			fmt.Printf("Failed to run script: %s\n", err.Error())
			panic(nil)
		}
	}

	if err != nil {
		fmt.Printf("Failed to get result: %s\n", err.Error())
		panic(nil)
	}

}

// BenchmarkFothFactorial allows running the foth benchmark
func BenchmarkFothFactorial(b *testing.B) {

	prg := `
: factorial recursive  dup 1 >  if  dup 1 -  factorial *  then  ;
100 factorial
`
	// Create
	f := fothEval.New()

	var err error

	// Run
	for i := 0; i < b.N; i++ {

		f.Reset()
		err = f.Eval(prg)

	}

	if err != nil {
		fmt.Printf("failed to run forth program: %s\n", err)
		panic(err)
	}

	// Get the result
	var res float64
	res, err = f.Stack.Pop()

	if err != nil {
		fmt.Printf("failed to get result: %s\n", err)
		panic(err)
	}

	if res != 93326215443944102188325606108575267240944254854960571509166910400407995064242937148632694030450512898042989296944474898258737204311236641477561877016501813248.000000 {
		fmt.Printf("Unexpected result:%f\n", res)
	}

}
