/*
Package pert helps you to work with estimated times to finish your project tasks
or activities.

PERT (Program Evaluation and Review Technique) is a statistical tool used in
project management that is designed to analyze and represent the tasks involved
in completing a given project.
*/
package pert

import (
	"errors"
	"math"
)

// Expect returns an expected time to accomplish a task or activity by given
// estimation times: optimistic time, realistic (most-likely) time, and
// pessimistic time.
// It returns an error if the estimates are not finite numbers satisfying
// 0 <= opt <= real <= pes.
func Expect(opt, real, pes float64) (float64, error) {
	if err := validateFinite(opt, real, pes); err != nil {
		return 0, err
	}
	if opt < 0 || opt > real || real > pes {
		return 0, errors.New("estimates must satisfy 0 <= optimistic <= realistic <= pessimistic")
	}

	weighted := opt + 4*real + pes
	if math.IsInf(weighted, 0) {
		// Scale the terms before adding when the weighted sum overflows.
		return opt/6 + (real/3)*2 + pes/6, nil
	}
	return weighted / 6, nil
}

// StdDev returns the PERT standard deviation for the given pessimistic and
// optimistic estimates.
// The larger results you get, the less confidence you have in your estimate.
// It returns an error if the estimates are not finite numbers satisfying
// 0 <= opt <= pes.
func StdDev(pes, opt float64) (float64, error) {
	if err := validateFinite(opt, pes); err != nil {
		return 0, err
	}
	if opt < 0 || opt > pes {
		return 0, errors.New("estimates must satisfy 0 <= optimistic <= pessimistic")
	}
	return (pes - opt) / 6, nil
}

func validateFinite(estimates ...float64) error {
	for _, value := range estimates {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			return errors.New("estimates must be finite numbers")
		}
	}
	return nil
}
