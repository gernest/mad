package console

import (
	"fmt"
	"time"

	"github.com/gernest/mad"
)

// Report pretty prints the spec results to stdout.
func Report(ts *mad.SpecResult) {
	printResult(ts, 0)
}

func printResult(ts *mad.SpecResult, level int) {
	fmt.Printf("%s%s: \n", ident(level), ts.Desc)
	for _, v := range ts.FailedExpectations {
		fmt.Printf("%s✖ %s :\n", ident(level+1), v.Desc)
		for _, msg := range v.Messages {
			fmt.Printf("%s-- %s \n", ident(level+2), msg)
		}
	}
	for _, v := range ts.PassedExpectations {
		fmt.Printf("%s✔ %s  %v\n", ident(level+1), v.Desc, v.Duration)
	}
	for _, v := range ts.Children {
		printResult(v, level+1)
	}
}

func ident(level int) string {
	s := ""
	for i := 0; i < level; i++ {
		s += "  "
	}
	return s
}

// calcStats calculates total failed and passed tests. This is recursive.
func calcStats(ts *mad.SpecResult) (int, int) {
	pass := len(ts.PassedExpectations)
	fail := len(ts.FailedExpectations)
	for _, v := range ts.Children {
		p, f := calcStats(v)
		pass += p
		fail += f
	}
	return pass, fail
}

// ResponseHandler implements mad.respHandler interface. This handles pretty
// printing spec results to stdout.
type ResponseHandler struct {
	passed   int
	failed   int
	Verbose  bool
	duration time.Duration
}

func New(verbose bool) *ResponseHandler {
	return &ResponseHandler{Verbose: verbose}
}

// Handle tracks the stats about the spec result and pretty prints the results to stdout.
func (r *ResponseHandler) Handle(ts *mad.SpecResult) {
	r.duration += ts.Duration
	pass, fail := calcStats(ts)
	r.passed += pass
	r.failed += fail
	if r.Verbose {
		Report(ts)
	} else {
		if fail > 0 {
			Report(ts)
		} else {
			fmt.Printf("%s✔ %s \n", ident(0), ts.Desc)
		}
	}
}

// Done prints the stats to stdout.
func (r *ResponseHandler) Done() {
	fmt.Printf(" Passed :%d Failed:%d in %s\n", r.passed, r.failed, r.duration)
}
