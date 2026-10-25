package pert_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/ink8bit/pert/pkg/pert"
)

func TestExpect(t *testing.T) {
	tests := []struct {
		name string
		opt  float64
		real float64
		pes  float64
		want float64
	}{
		{name: "distinct estimates", opt: 1, real: 3, pes: 13, want: 26.0 / 6},
		{name: "equal estimates", opt: 3, real: 3, pes: 3, want: 3},
		{name: "zero estimates", opt: 0, real: 0, pes: 0, want: 0},
		{name: "zero optimistic estimate", opt: 0, real: 1, pes: 6, want: 10.0 / 6},
		{name: "fractional estimates", opt: 0.5, real: 1.5, pes: 2.5, want: 1.5},
		{name: "large equal estimates", opt: 1e308, real: 1e308, pes: 1e308, want: 1e308},
		{name: "maximum equal estimates", opt: math.MaxFloat64, real: math.MaxFloat64, pes: math.MaxFloat64, want: math.MaxFloat64},
		{name: "multiplication overflow", opt: 0, real: math.MaxFloat64 / 2, pes: math.MaxFloat64, want: math.MaxFloat64 / 2},
		{name: "addition overflow", opt: math.MaxFloat64 / 4, real: math.MaxFloat64 / 4, pes: math.MaxFloat64 / 2, want: math.MaxFloat64 * 7 / 24},
		{name: "large distinct estimates", opt: math.MaxFloat64 / 4, real: math.MaxFloat64 / 2, pes: math.MaxFloat64, want: math.MaxFloat64 * 13 / 24},
		{name: "maximum pessimistic estimate", opt: 0, real: 0, pes: math.MaxFloat64, want: math.MaxFloat64 / 6},
		{name: "smallest equal estimates", opt: math.SmallestNonzeroFloat64, real: math.SmallestNonzeroFloat64, pes: math.SmallestNonzeroFloat64, want: math.SmallestNonzeroFloat64},
		{name: "smallest distinct estimates", opt: math.SmallestNonzeroFloat64, real: 2 * math.SmallestNonzeroFloat64, pes: 3 * math.SmallestNonzeroFloat64, want: 2 * math.SmallestNonzeroFloat64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pert.Expect(tt.opt, tt.real, tt.pes)
			if err != nil {
				t.Fatalf("Expect(%g, %g, %g) returned an error: %v", tt.opt, tt.real, tt.pes, err)
			}
			if math.IsNaN(got) || math.IsInf(got, 0) || math.Abs(got-tt.want) > 1e-15*math.Abs(tt.want) {
				t.Errorf("Expect(%g, %g, %g) = %g, want %g", tt.opt, tt.real, tt.pes, got, tt.want)
			}
			if got < tt.opt || got > tt.pes {
				t.Errorf("expected value %g is outside [%g, %g]", got, tt.opt, tt.pes)
			}
		})
	}
}

func TestStdDev(t *testing.T) {
	tests := []struct {
		name string
		opt  float64
		pes  float64
		want float64
	}{
		{name: "distinct estimates", opt: 1, pes: 13, want: 2},
		{name: "equal estimates", opt: 3, pes: 3, want: 0},
		{name: "zero estimates", opt: 0, pes: 0, want: 0},
		{name: "zero optimistic estimate", opt: 0, pes: 6, want: 1},
		{name: "fractional estimates", opt: 0.5, pes: 2, want: 0.25},
		{name: "maximum pessimistic estimate", opt: 0, pes: math.MaxFloat64, want: math.MaxFloat64 / 6},
		{name: "maximum equal estimates", opt: math.MaxFloat64, pes: math.MaxFloat64, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pert.StdDev(tt.opt, tt.pes)
			if err != nil {
				t.Fatalf("StdDev(%g, %g) returned an error: %v", tt.opt, tt.pes, err)
			}
			if got != tt.want {
				t.Errorf("StdDev(%g, %g) = %g, want %g", tt.opt, tt.pes, got, tt.want)
			}
		})
	}
}

func TestExpectInvalidEstimates(t *testing.T) {
	const ordered = "estimates must satisfy 0 <= optimistic <= realistic <= pessimistic"
	type testCase struct {
		name      string
		estimates [3]float64
		want      string
	}
	tests := []testCase{
		{name: "negative optimistic", estimates: [3]float64{-1, 3, 13}, want: ordered},
		{name: "negative realistic", estimates: [3]float64{1, -3, 13}, want: ordered},
		{name: "negative pessimistic", estimates: [3]float64{1, 3, -13}, want: ordered},
		{name: "all negative", estimates: [3]float64{-3, -2, -1}, want: ordered},
		{name: "realistic below optimistic", estimates: [3]float64{5, 3, 13}, want: ordered},
		{name: "realistic above pessimistic", estimates: [3]float64{1, 15, 13}, want: ordered},
		{name: "descending estimates", estimates: [3]float64{13, 3, 1}, want: ordered},
	}
	for i, name := range []string{"optimistic", "realistic", "pessimistic"} {
		for _, value := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
			estimates := [3]float64{1, 3, 13}
			estimates[i] = value
			tests = append(tests, testCase{
				name: name + " " + fmt.Sprint(value), estimates: estimates, want: "estimates must be finite numbers",
			})
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pert.Expect(tt.estimates[0], tt.estimates[1], tt.estimates[2])
			if err == nil || err.Error() != tt.want {
				t.Errorf("Expect(%g, %g, %g) error = %v, want %q", tt.estimates[0], tt.estimates[1], tt.estimates[2], err, tt.want)
			}
			if got != 0 {
				t.Errorf("Expect(%g, %g, %g) = %g, want 0 for invalid estimates", tt.estimates[0], tt.estimates[1], tt.estimates[2], got)
			}
		})
	}
}

func TestStdDevInvalidEstimates(t *testing.T) {
	const ordered = "estimates must satisfy 0 <= optimistic <= pessimistic"
	type testCase struct {
		name string
		opt  float64
		pes  float64
		want string
	}
	tests := []testCase{
		{name: "negative optimistic", opt: -1, pes: 13, want: ordered},
		{name: "negative pessimistic", opt: 1, pes: -13, want: ordered},
		{name: "both negative", opt: -3, pes: -1, want: ordered},
		{name: "pessimistic below optimistic", opt: 13, pes: 1, want: ordered},
	}
	for _, value := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		tests = append(tests,
			testCase{name: "optimistic " + fmt.Sprint(value), opt: value, pes: 13, want: "estimates must be finite numbers"},
			testCase{name: "pessimistic " + fmt.Sprint(value), opt: 1, pes: value, want: "estimates must be finite numbers"},
		)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pert.StdDev(tt.opt, tt.pes)
			if err == nil || err.Error() != tt.want {
				t.Errorf("StdDev(%g, %g) error = %v, want %q", tt.opt, tt.pes, err, tt.want)
			}
			if got != 0 {
				t.Errorf("StdDev(%g, %g) = %g, want 0 for invalid estimates", tt.opt, tt.pes, got)
			}
		})
	}
}
