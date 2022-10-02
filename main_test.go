package main

import (
	"fmt"
	"testing"

	// critical
	"github.com/skx/critical/interpreter"
	criticalSTDLIB	"github.com/skx/critical/stdlib"

	// yal
	"github.com/skx/yal/builtins"
	"github.com/skx/yal/env"
	"github.com/skx/yal/eval"
	"github.com/skx/yal/primitive"
	"github.com/skx/yal/stdlib"
)

// The YAL interpreter we're going to execute with.
var yalInter *eval.Eval

// The environment contains the primitives the interpreter uses.
var yalEnv *env.Environment

var tclInter *interpreter.Interpreter


// Create the various interpreters, and prime them for execution.
//
// Only do this once, at startup.
func init() {

	initYAL()

	initTCL()
}

// initTCL sets up a default critical/TCL interpreter.
func initTCL() {

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
	stdlib := criticalSTDLIB.Contents()

	var err error
	tclInter, err = interpreter.New(string(stdlib) + "\n" + input)
	if err != nil {
		fmt.Printf("initTCL; Error creating interpreter %s\n", err)
		panic(err)
		return
	}
}

// initYAL sets up a default YAL interpreter, ready to run the benchmark
func initYAL() {
	// Create a new environment
	yalEnv = env.New()

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
	yalInter = eval.New(src)
}

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
	var out primitive.Primitive

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
	var out any
	var err error
	for i := 0; i < b.N; i++ {

		// Run 100!
		out, err = tclInter.Evaluate()

	}
	if err != nil && err != interpreter.ErrReturn {
		fmt.Printf("Error running program:%s\n", err)
		return
	}
	if out != "93326215443944102188325606108575267240944254854960571509166910400407995064242937148632694030450512898042989296944474898258737204311236641477561877016501813248.000000" {
		fmt.Printf("Unexpected result:%s\n", out)
	}

}
