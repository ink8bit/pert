package main

import (
	"flag"
	"fmt"
	"io"
	"math"
	"os"

	"github.com/ink8bit/pert/pkg/pert"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("pert", flag.ContinueOnError)
	flags.SetOutput(stderr)

	var opt float64
	flags.Float64Var(&opt, "o", 0, "optimistic value (required)")

	var real float64
	flags.Float64Var(&real, "r", 0, "realistic value (required)")

	var pes float64
	flags.Float64Var(&pes, "p", 0, "pessimistic value (required)")

	var name string
	flags.StringVar(&name, "n", "", "task name")

	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}

	provided := make(map[string]bool)
	flags.Visit(func(f *flag.Flag) {
		provided[f.Name] = true
	})
	if !provided["o"] || !provided["r"] || !provided["p"] {
		fmt.Fprintln(stderr, "all estimates are required: -o, -r, -p")
		return 2
	}

	for _, value := range []float64{opt, real, pes} {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			fmt.Fprintln(stderr, "estimates must be finite numbers")
			return 2
		}
	}
	if opt < 0 || opt > real || real > pes {
		fmt.Fprintln(stderr, "estimates must satisfy 0 <= optimistic <= realistic <= pessimistic")
		return 2
	}

	e := pert.Expect(opt, real, pes)
	v := pert.Variance(pes, opt)

	fmt.Fprintln(stdout, output(name, e, v))
	return 0
}

func output(label string, exp, dev float64) string {
	fmtStr := fmt.Sprintf("- Expected:\t%.2f\n- Deviation:\t%.2f", exp, dev)
	if label != "" {
		return fmt.Sprintf("[%s]\n", label) + fmtStr
	}
	return fmtStr
}
