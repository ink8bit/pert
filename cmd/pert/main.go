package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ink8bit/pert/internal/version"
	"github.com/ink8bit/pert/pkg/pert"
)

func main() {
	var opt float64
	flag.Float64Var(&opt, "o", 0, "optimistic value")

	var real float64
	flag.Float64Var(&real, "r", 0, "realistic value")

	var pes float64
	flag.Float64Var(&pes, "p", 0, "pessimistic value")

	var name string
	flag.StringVar(&name, "n", "", "task name")

	var ver bool
	flag.BoolVar(&ver, "v", false, "print version")

	flag.Parse()

	if ver {
		fmt.Println(version.Print())
		os.Exit(0)
	}

	args := args{
		Opt:  opt,
		Real: real,
		Pes:  pes,
	}

	res, err := run(args)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(output(name, res.Expected, res.Variance))
}

type args struct {
	Opt  float64
	Real float64
	Pes  float64
}

type result struct {
	Expected float64
	Variance float64
}

func run(args args) (*result, error) {
	if args.Opt == 0 && args.Real == 0 && args.Pes == 0 {
		return nil, errors.New("No values provided")
	}
	e := pert.Expect(args.Opt, args.Real, args.Pes)
	v := pert.Variance(args.Pes, args.Opt)
	r := &result{
		Expected: e,
		Variance: v,
	}
	return r, nil
}

func output(label string, exp, dev float64) string {
	fmtStr := fmt.Sprintf("- Expected:\t%.2f\n- Deviation:\t%.2f", exp, dev)
	if label != "" {
		return fmt.Sprintf("[%s]\n", label) + fmtStr
	}
	return fmtStr
}
