package main

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestRunValidEstimates(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "named task",
			args: []string{"-n", "Example", "-o", "1", "-r", "3", "-p", "13"},
			want: "[Example]\n- Expected:\t4.33\n- Deviation:\t2.00\n",
		},
		{
			name: "explicit zeros",
			args: []string{"-o", "0", "-r", "0", "-p", "0"},
			want: "- Expected:\t0.00\n- Deviation:\t0.00\n",
		},
		{
			name: "zero optimistic estimate",
			args: []string{"-o", "0", "-r", "1", "-p", "6"},
			want: "- Expected:\t1.67\n- Deviation:\t1.00\n",
		},
		{
			name: "equal estimates",
			args: []string{"-o", "3", "-r", "3", "-p", "3"},
			want: "- Expected:\t3.00\n- Deviation:\t0.00\n",
		},
		{
			name: "realistic equals optimistic",
			args: []string{"-o", "1", "-r", "1", "-p", "7"},
			want: "- Expected:\t2.00\n- Deviation:\t1.00\n",
		},
		{
			name: "realistic equals pessimistic",
			args: []string{"-o", "1", "-r", "7", "-p", "7"},
			want: "- Expected:\t6.00\n- Deviation:\t1.00\n",
		},
		{
			name: "fractional estimates",
			args: []string{"-o", "0.5", "-r", "1.5", "-p", "2.5"},
			want: "- Expected:\t1.50\n- Deviation:\t0.33\n",
		},
		{
			name: "large finite estimates",
			args: []string{"-o", "1e308", "-r", "1e308", "-p", "1e308"},
			want: fmt.Sprintf("- Expected:\t%.2f\n- Deviation:\t0.00\n", 1e308),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if code := run(tt.args, &stdout, &stderr); code != 0 {
				t.Fatalf("exit code = %d, want 0; stderr: %s", code, &stderr)
			}
			if got := stdout.String(); got != tt.want {
				t.Errorf("stdout = %q, want %q", got, tt.want)
			}
			if stderr.Len() != 0 {
				t.Errorf("unexpected stderr: %s", &stderr)
			}
		})
	}
}

func TestRunInvalidEstimates(t *testing.T) {
	const required = "all estimates are required: -o, -r, -p"
	const ordered = "0 <= optimistic <= realistic <= pessimistic"
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "no estimates", want: required},
		{name: "only optimistic", args: []string{"-o", "12"}, want: required},
		{name: "only realistic", args: []string{"-r", "3"}, want: required},
		{name: "only pessimistic", args: []string{"-p", "13"}, want: required},
		{name: "missing optimistic", args: []string{"-r", "3", "-p", "13"}, want: required},
		{name: "missing realistic", args: []string{"-o", "1", "-p", "13"}, want: required},
		{name: "missing pessimistic", args: []string{"-o", "1", "-r", "3"}, want: required},
		{name: "missing zero estimate", args: []string{"-o", "0", "-r", "0"}, want: required},
		{name: "negative optimistic", args: []string{"-o=-1", "-r", "3", "-p", "13"}, want: ordered},
		{name: "negative realistic", args: []string{"-o", "1", "-r=-3", "-p", "13"}, want: ordered},
		{name: "negative pessimistic", args: []string{"-o", "1", "-r", "3", "-p=-13"}, want: ordered},
		{name: "all negative", args: []string{"-o=-3", "-r=-2", "-p=-1"}, want: ordered},
		{name: "realistic below optimistic", args: []string{"-o", "5", "-r", "3", "-p", "13"}, want: ordered},
		{name: "realistic above pessimistic", args: []string{"-o", "1", "-r", "15", "-p", "13"}, want: ordered},
	}

	for i, flagName := range []string{"-o", "-r", "-p"} {
		for _, value := range []string{"NaN", "+Inf", "-Inf"} {
			args := []string{"-o", "1", "-r", "3", "-p", "13"}
			args[2*i+1] = value
			tests = append(tests, struct {
				name string
				args []string
				want string
			}{name: flagName + " " + value, args: args, want: "estimates must be finite numbers"})
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if code := run(tt.args, &stdout, &stderr); code != 2 {
				t.Errorf("exit code = %d, want 2", code)
			}
			if stdout.Len() != 0 {
				t.Errorf("unexpected stdout: %s", &stdout)
			}
			if got := stderr.String(); !strings.Contains(got, tt.want) {
				t.Errorf("stderr = %q, want it to contain %q", got, tt.want)
			}
		})
	}
}

func TestRunUnexpectedArguments(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "trailing argument",
			args: []string{"-o", "1", "-r", "3", "-p", "13", "unexpected"},
			want: "unexpected arguments: [\"unexpected\"]\n",
		},
		{
			name: "argument before flags",
			args: []string{"unexpected", "-o", "1", "-r", "3", "-p", "13"},
			want: "unexpected arguments: [\"unexpected\" \"-o\" \"1\" \"-r\" \"3\" \"-p\" \"13\"]\n",
		},
		{
			name: "argument between flags",
			args: []string{"-o", "1", "unexpected", "-r", "3", "-p", "13"},
			want: "unexpected arguments: [\"unexpected\" \"-r\" \"3\" \"-p\" \"13\"]\n",
		},
		{
			name: "flags after trailing argument",
			args: []string{"-o", "1", "-r", "3", "-p", "13", "unexpected", "-r", "100"},
			want: "unexpected arguments: [\"unexpected\" \"-r\" \"100\"]\n",
		},
		{
			name: "arguments after separator",
			args: []string{"-o", "1", "-r", "3", "-p", "13", "--", "-r", "100"},
			want: "unexpected arguments: [\"-r\" \"100\"]\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if code := run(tt.args, &stdout, &stderr); code != 2 {
				t.Errorf("exit code = %d, want 2", code)
			}
			if stdout.Len() != 0 {
				t.Errorf("unexpected stdout: %s", &stdout)
			}
			if got := stderr.String(); got != tt.want {
				t.Errorf("stderr = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRunOutputError(t *testing.T) {
	writeErr := errors.New("output unavailable")
	for _, tt := range []struct {
		name    string
		written int
	}{
		{name: "complete write failure", written: 0},
		{name: "partial write failure", written: 5},
	} {
		t.Run(tt.name, func(t *testing.T) {
			stdout := failingWriter{written: tt.written, err: writeErr}
			var stderr bytes.Buffer
			if code := run([]string{"-o", "1", "-r", "3", "-p", "13"}, stdout, &stderr); code != 1 {
				t.Errorf("exit code = %d, want 1", code)
			}
			if got, want := stderr.String(), "failed to write output: output unavailable\n"; got != want {
				t.Errorf("stderr = %q, want %q", got, want)
			}
		})
	}
}

type failingWriter struct {
	written int
	err     error
}

func (w failingWriter) Write([]byte) (int, error) {
	return w.written, w.err
}

func TestRunHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := run([]string{"-h"}, &stdout, &stderr); code != 0 {
		t.Errorf("exit code = %d, want 0", code)
	}
	if stdout.Len() != 0 {
		t.Errorf("unexpected stdout: %s", &stdout)
	}
	if got := stderr.String(); !strings.Contains(got, "Usage of pert:") {
		t.Errorf("stderr = %q, want help text", got)
	}
}
