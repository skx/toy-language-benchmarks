package main

import (
	"fmt"
	"testing"

	// basic
	basicEVAL "github.com/skx/gobasic/eval"
	basicOBJECT "github.com/skx/gobasic/object"
	basicTOKENIZER "github.com/skx/gobasic/tokenizer"

	// critical
	"github.com/skx/critical/interpreter"
	criticalSTDLIB "github.com/skx/critical/stdlib"

	// evalfilter
	"github.com/skx/evalfilter/v2"
	evalObject "github.com/skx/evalfilter/v2/object"

	// foth
	fothEval "github.com/skx/foth/foth/eval"

	// monkey
	monkeyEval "github.com/skx/monkey/evaluator"
	monkeyLexer "github.com/skx/monkey/lexer"
	monkeyObject "github.com/skx/monkey/object"
	monkeyParser "github.com/skx/monkey/parser"

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
	var out int64
	for i := 0; i < b.N; i++ {
		out = fact(12)
	}
	if out != 479001600 {
		b.Fatalf("unexpected result %d", out)
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

(fact 12)
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

		// Run 12!
		out = yalInter.Evaluate(yalEnv)
	}

	// Did we get an error?  Then show it.
	if _, ok := out.(primitive.Error); ok {
		fmt.Printf("Error running: %v\n", out)
	}

	if out.ToString() != "479001600" {
		b.Fatalf("unexpected result %s", out.ToString())
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
fact 12
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

		// Run 12!
		out, err = tclInter.Evaluate()

	}

	// Ensure the result was "ok".
	if err != nil && err != interpreter.ErrReturn {
		fmt.Printf("Error running program:%s\n", err)
		return
	}
	if out != "479001600" {
		b.Fatalf("Unexpected result:%s\n", out)
	}
}

// BenchmarkEvalFilterFactorial allows running the evalfilter benchmark.
func BenchmarkEvalFilterFactorial(b *testing.B) {

	prg := `
function fact( n ) {
  if ( n <= 1 ) { return 1; }
  return ( n * fact( n - 1 ) );
}
return fact(12);
`
	eval := evalfilter.New(string(prg))

	err := eval.Prepare([]byte{})
	var out evalObject.Object

	if err != nil {
		fmt.Printf("Error compiling:%s\n", err.Error())
		panic(err)
	}

	for i := 0; i < b.N; i++ {

		out, err = eval.Execute(true)
		if err != nil {
			fmt.Printf("Failed to run script: %s\n", err.Error())
			panic(nil)
		}
	}

	if err != nil {
		fmt.Printf("Failed to get result: %s\n", err.Error())
		panic(nil)
	}

	if out.Inspect() != "479001600" {
		b.Fatalf("unexpected results %s", out)
	}
}

// BenchmarkFothFactorial allows running the forth benchmark.
func BenchmarkFothFactorial(b *testing.B) {

	prg := `
: factorial recursive  dup 1 >  if  dup 1 -  factorial *  then  ;
12 factorial
`

	var err error
	var f *fothEval.Eval

	// Run
	for i := 0; i < b.N; i++ {

		b.StopTimer()

		// Create
		f = fothEval.New()

		//
		// BUG: "f.Reset()" doesn't do enough???
		//

		b.StartTimer()

		// Execute
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

	if res != 479001600 {
		b.Fatalf("Unexpected result %f", res)
	}
}

// BenchmarkMonkeyFactorial allows running the monkey benchmark.
func BenchmarkMonkeyFactorial(b *testing.B) {

	prg := `
function fact( n ) {
  if ( n <= 1 ) { return 1; }
  return ( n * fact( n - 1 ) );
}
return fact(12);
`

	env := monkeyObject.NewEnvironment()
	l := monkeyLexer.New(prg)
	p := monkeyParser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		for _, msg := range p.Errors() {
			fmt.Printf("\t%s\n", msg)
		}
		panic("failed to parse program")
	}

	var out monkeyObject.Object

	// Run
	for i := 0; i < b.N; i++ {
		out = monkeyEval.Eval(program, env)
	}

	if out.Inspect() != "479001600" {
		b.Fatalf("Unexpected result %s", out.Inspect())
	}
}

// BenchmarkBASICFactorial allows running the BASIC benchmark.
func BenchmarkBASICFactorial(b *testing.B) {
	prg := `
10 LET F = 0
20 LET N = 12
20 GOSUB 200
30 END


REM
REM This routine calculates a factorial, recursively.
REM
REM Before calling set "F=0", and store the number you wish to calculate in N
REM
REM So to calculate "5!" run
REM
REM  LET F = 0
REM  LET N = 5
REM  GOSUB 200
REM  PRINT "Result: ", f, "\n"
REM
REM Here we suffer from one-statement per line restrictions, as made obvious
REM in the duplicated conditionals.
REM
200 IF N<0 THEN F=-1
210 IF N<0 THEN RETURN
220 IF N<2 THEN F=1
230 IF N<2 THEN RETURN
240 LET N = N - 1
250 GOSUB 200
260 LET N = N + 1
270 LET F = F * N
280 RETURN
`

	var err error
	var result basicOBJECT.Object

	for i := 0; i < b.N; i++ {

		// There's no reset of state in our BASIC interpreter
		//
		// So we need to reparse from scratch each time, but that
		// shouldn't count towards the benchmark time.  So we stop/start
		// the timer around that setup.
		//
		b.StopTimer()

		t := basicTOKENIZER.New(string(prg))
		e, err2 := basicEVAL.New(t)
		if err2 != nil {
			fmt.Printf("Failed to parse program: %s\n", err)
			panic(err)
		}

		b.StartTimer()

		err = e.Run()

		result = e.GetVariable("F")
	}

	if err != nil {
		fmt.Printf("failed to run program: %s\n", err)
		panic(err)
	}

	// Ensure the result was a number
	num, ok := result.(*basicOBJECT.NumberObject)
	if !ok {
		b.Fatalf("didn't get a number result, got %s", result)
	}

	// Of the correct value.
	if num.Value != 479001600 {
		b.Fatalf("Unexpected result %s", result.String())
	}
}
