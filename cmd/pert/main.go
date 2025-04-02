package main

import (
	"flag"
	"fmt"

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

	flag.Parse()

	if opt == 0 && real == 0 && pes == 0 {
		fmt.Println("No value provided")
		return
	}

	e := pert.Expect(opt, real, pes)
	v := pert.Variance(pes, opt)

	fmt.Println(output(name, e, v))
}

func output(label string, exp, dev float64) string {
	fmtStr := fmt.Sprintf("- Expected:\t%.2f\n- Deviation:\t%.2f", exp, dev)
	if label != "" {
		return fmt.Sprintf("[%s]\n", label) + fmtStr
	}
	return fmtStr
}
