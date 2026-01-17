package diagnostics

import (
	"fmt"
	"io"
)

// Report contains all diagnostic check results.
type Report struct {
	Checks []Check
}

// NewReport creates a new empty report.
func NewReport() *Report {
	return &Report{
		Checks: []Check{},
	}
}

// Add adds a check to the report.
func (r *Report) Add(c Check) {
	r.Checks = append(r.Checks, c)
}

// AddAll adds multiple checks to the report.
func (r *Report) AddAll(checks []Check) {
	r.Checks = append(r.Checks, checks...)
}

// HasFailures returns true if any check failed.
func (r *Report) HasFailures() bool {
	for _, c := range r.Checks {
		if c.Status == StatusFail {
			return true
		}
	}
	return false
}

// HasWarnings returns true if any check has warnings.
func (r *Report) HasWarnings() bool {
	for _, c := range r.Checks {
		if c.Status == StatusWarn {
			return true
		}
	}
	return false
}

// Print outputs the report to the given writer.
func (r *Report) Print(w io.Writer) {
	for _, c := range r.Checks {
		icon := statusIcon(c.Status)
		fmt.Fprintf(w, "[%s] %s: %s\n", icon, c.Name, c.Message)
		if c.Details != "" && c.Status != StatusOK {
			fmt.Fprintf(w, "    %s\n", c.Details)
		}
	}
}

// PrintVerbose outputs the report with full details.
func (r *Report) PrintVerbose(w io.Writer) {
	for _, c := range r.Checks {
		icon := statusIcon(c.Status)
		fmt.Fprintf(w, "[%s] %s: %s\n", icon, c.Name, c.Message)
		if c.Details != "" {
			fmt.Fprintf(w, "    %s\n", c.Details)
		}
	}
}

func statusIcon(s Status) string {
	switch s {
	case StatusOK:
		return "\u2713" // ✓
	case StatusWarn:
		return "!"
	case StatusFail:
		return "\u2717" // ✗
	default:
		return "?"
	}
}
