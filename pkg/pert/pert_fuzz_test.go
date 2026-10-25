package pert_test

import (
	"math"
	"math/big"
	"slices"
	"testing"

	"github.com/ink8bit/pert/pkg/pert"
)

func FuzzExpect(f *testing.F) {
	for _, seed := range [][3]float64{
		{0, 0, 0},
		{1, 3, 13},
		{0.5, 1.5, 2.5},
		{1e308, 1e308, 1e308},
		{math.MaxFloat64, math.MaxFloat64, math.MaxFloat64},
		{0, math.MaxFloat64 / 2, math.MaxFloat64},
		{math.MaxFloat64 / 4, math.MaxFloat64 / 4, math.MaxFloat64 / 2},
		{0, 0, math.MaxFloat64},
		{math.MaxFloat64 / 6, math.MaxFloat64 / 6, math.MaxFloat64 / 6},
		{0, math.Nextafter(math.MaxFloat64/4, math.Inf(1)), math.MaxFloat64},
		{math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64},
		{0, math.SmallestNonzeroFloat64, 2 * math.SmallestNonzeroFloat64},
	} {
		f.Add(seed[0], seed[1], seed[2])
	}

	f.Fuzz(func(t *testing.T, a, b, c float64) {
		estimates := []float64{a, b, c}
		for _, value := range estimates {
			if value < 0 || math.IsNaN(value) || math.IsInf(value, 0) {
				t.Skip()
			}
		}
		slices.Sort(estimates)
		opt, real, pes := estimates[0], estimates[1], estimates[2]

		got, err := pert.Expect(opt, real, pes)
		if err != nil {
			t.Fatalf("Expect(%g, %g, %g) returned an error: %v", opt, real, pes, err)
		}
		if math.IsNaN(got) || math.IsInf(got, 0) {
			t.Fatalf("Expect(%g, %g, %g) = %g, want a finite result", opt, real, pes, got)
		}

		// Rational arithmetic provides an exact reference without float64 overflow.
		expected := new(big.Rat).Mul(new(big.Rat).SetFloat64(real), big.NewRat(4, 1))
		expected.Add(expected, new(big.Rat).SetFloat64(opt))
		expected.Add(expected, new(big.Rat).SetFloat64(pes))
		expected.Quo(expected, big.NewRat(6, 1))
		want, _ := expected.Float64()

		// Allow four ULPs of rounding error. The downward spacing stays finite at
		// MaxFloat64 and provides an absolute tolerance for subnormal results.
		tolerance := 4 * (want - math.Nextafter(want, math.Inf(-1)))
		if math.Abs(got-want) > tolerance {
			t.Errorf("Expect(%g, %g, %g) = %g, want %g within %g", opt, real, pes, got, want, tolerance)
		}
		if (got < opt && opt-got > tolerance) || (got > pes && got-pes > tolerance) {
			t.Errorf("expected value %g is outside [%g, %g] beyond tolerance %g", got, opt, pes, tolerance)
		}
	})
}
