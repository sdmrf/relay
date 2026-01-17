package app

import "github.com/sdmrf/relay/internal/plan"

// Executor executes a plan.
// Implementations: FSExecutor (real), dry-run, logging wrappers.
type Executor interface {
	Execute(plan.Plan) error
}
