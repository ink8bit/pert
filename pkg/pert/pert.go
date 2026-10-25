/*
Package pert helps you to work with estimated times to finish your project tasks
or activities.

PERT (Program Evaluation and Review Technique) is a statistical tool used in
project management that is designed to analyze and represent the tasks involved
in completing a given project.
*/
package pert

// Expect returns an expected time to accomplish a task or activity by given
// estimation times: optimistic time, realistic (most-likely) time, and
// pessimistic time.
func Expect(opt, real, pes float64) float64 {
	return (opt + 4*real + pes) / 6
}

// StdDev returns the PERT standard deviation for the given pessimistic and
// optimistic estimates.
// The larger results you get, the less confidence you have in your estimate.
func StdDev(pes, opt float64) float64 {
	return (pes - opt) / 6
}
