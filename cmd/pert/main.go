package main

import (
	_ "embed"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ink8bit/pert/pkg/pert"
)

//go:embed version.txt
var version string

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("pert", flag.ContinueOnError)
	flags.SetOutput(stderr)

	var showVersion bool
	flags.BoolVar(&showVersion, "version", false, "show version and exit")

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

	if flags.NArg() != 0 {
		fmt.Fprintf(stderr, "unexpected arguments: %q\n", flags.Args())
		return 2
	}

	if showVersion {
		return writeOutput(stdout, stderr, strings.TrimSpace(version))
	}

	provided := make(map[string]bool)
	flags.Visit(func(f *flag.Flag) {
		provided[f.Name] = true
	})
	if !provided["o"] || !provided["r"] || !provided["p"] {
		fmt.Fprintln(stderr, "all estimates are required: -o, -r, -p")
		return 2
	}

	e, err := pert.Expect(opt, real, pes)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	dev, err := pert.StdDev(opt, pes)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}

	return writeOutput(stdout, stderr, output(name, e, dev))
}

func writeOutput(stdout, stderr io.Writer, value string) int {
	if _, err := fmt.Fprintln(stdout, value); err != nil {
		fmt.Fprintf(stderr, "failed to write output: %v\n", err)
		return 1
	}
	return 0
}

func output(label string, exp, dev float64) string {
	fmtStr := fmt.Sprintf("- Expected:\t%.2f\n- Deviation:\t%.2f", exp, dev)
	if label != "" {
		return fmt.Sprintf("[%s]\n", label) + fmtStr
	}
	return fmtStr
}
