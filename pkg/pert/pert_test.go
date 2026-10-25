package pert_test

import (
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
			got := pert.Expect(tt.opt, tt.real, tt.pes)
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
		pes  float64
		opt  float64
		want float64
	}{
		{name: "distinct estimates", pes: 13, opt: 1, want: 2},
		{name: "equal estimates", pes: 3, opt: 3, want: 0},
		{name: "zero estimates", pes: 0, opt: 0, want: 0},
		{name: "zero optimistic estimate", pes: 6, opt: 0, want: 1},
		{name: "fractional estimates", pes: 2, opt: 0.5, want: 0.25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pert.StdDev(tt.pes, tt.opt); got != tt.want {
				t.Errorf("StdDev(%g, %g) = %g, want %g", tt.pes, tt.opt, got, tt.want)
			}
		})
	}
}
