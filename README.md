# pert CLI

[![Go Reference](https://pkg.go.dev/badge/github.com/ink8bit/pert.svg)](https://pkg.go.dev/github.com/ink8bit/pert)

PERT (Program Evaluation and Review Technique) is a statistical tool used in
project management that is designed to analyze and represent the tasks involved
in completing a given project.

## Install

CLI tool:

```console
go install github.com/ink8bit/pert/cmd/pert@latest
```

Package:

```console
go get github.com/ink8bit/pert
```

## Usage

All three estimates (`-o`, `-r`, and `-p`) are required. They must be finite
numbers satisfying `0 <= optimistic <= realistic <= pessimistic`. Zero estimates
are allowed when explicitly supplied.

```console
Usage of pert:
  -n string
        task name
  -o float
        optimistic value (required)
  -p float
        pessimistic value (required)
  -r float
        realistic value (required)
```

## Library

`pert.Expect(opt, real, pes)` and `pert.StdDev(pes, opt)` return `(float64, error)`.
Both reject non-finite or negative estimates. `Expect` requires
`optimistic <= realistic <= pessimistic`; `StdDev` requires
`optimistic <= pessimistic`. Equal and zero estimates are valid.

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
    deviation, err := pert.StdDev(13, 1)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Expected: %.2f, deviation: %.2f\n", expected, deviation)
}
```
