package pert_test

import (
	"testing"

	"github.com/ink8bit/pert/pkg/pert"
)

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
