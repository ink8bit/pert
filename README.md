# pert

[![Go Reference](https://pkg.go.dev/badge/github.com/ink8bit/pert/pkg/pert.svg)](https://pkg.go.dev/github.com/ink8bit/pert/pkg/pert)

A Go library and command-line tool for estimating a task's expected duration and
standard deviation using PERT (Program Evaluation and Review Technique).

## Install

Requires Go 1.24 or later.

CLI tool:

```console
go install github.com/ink8bit/pert/cmd/pert@latest
```

Package:

```console
go get github.com/ink8bit/pert/pkg/pert
```

## Usage

For a task with optimistic, realistic, and pessimistic estimates of 1, 3, and 13
hours:

```console
$ pert -n "Example task" -o 1 -r 3 -p 13
[Example task]
- Expected:     4.33
- Deviation:    2.00
```

The expected duration is approximately 4.33 hours, with an estimated standard
deviation of 2 hours. Use the same unit for all estimates, such as hours or days;
both results use that unit. The tool does not perform unit conversion.

| Flag | Description | Required |
| --- | --- | --- |
| `-o` | Optimistic duration | Yes |
| `-r` | Realistic (most likely) duration | Yes |
| `-p` | Pessimistic duration | Yes |
| `-n` | Task name, shown above the results | No |
| `-h` | Show help | No |

All three estimates (`-o`, `-r`, and `-p`) are required. They must be finite
numbers satisfying `0 <= optimistic <= realistic <= pessimistic`. Zero estimates
are allowed when explicitly supplied; negative values, `NaN`, and infinities are
rejected. Positional arguments are rejected.

The CLI prints results to two decimal places, so very small positive results can
appear as `0.00`. Use the library when you need the unformatted `float64` values.

Results are written to standard output; help and errors are written to standard
error. Exit codes are:

| Code | Meaning |
| --- | --- |
| `0` | Successful calculation or help request |
| `1` | Failed to write results |
| `2` | Invalid flags, missing or invalid estimates, or unexpected arguments |

## Calculations

The library uses these formulas:

```text
expected  = (optimistic + 4 * realistic + pessimistic) / 6
deviation = (pessimistic - optimistic) / 6
```

The realistic estimate is the most likely duration; it receives four times the
weight of either extreme in the expected duration. Deviation describes the
estimated spread in the same unit as the duration. Its magnitude depends on the
unit and task scale, so it is not a confidence percentage.

## Library

`pert.Expect(opt, real, pes)` and `pert.StdDev(opt, pes)` return `(float64, error)`.
Both reject non-finite or negative estimates. `Expect` requires
`optimistic <= realistic <= pessimistic`; `StdDev` requires
`optimistic <= pessimistic`. Equal and zero estimates are valid.

Both functions take estimates from optimistic to pessimistic. `StdDev` previously
accepted the reverse order; existing callers must swap its arguments.

```go
package main

import (
	"fmt"
	"log"

	"github.com/ink8bit/pert/pkg/pert"
)

func main() {
	expected, err := pert.Expect(1, 3, 13)
	if err != nil {
		log.Fatal(err)
	}
	deviation, err := pert.StdDev(1, 13)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Expected: %.2f, deviation: %.2f\n", expected, deviation)
}
```

Output:

```text
Expected: 4.33, deviation: 2.00
```

## Development

Run tests, static checks, and a build from the repository root:

```sh
go test ./...
go vet ./...
go build ./...
```
