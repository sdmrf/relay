package app

import (
	"context"

	"github.com/sdmrf/relay/internal/plan"
)

// Executor executes a plan.
// Implementations: FSExecutor (real), dry-run, logging wrappers.
type Executor interface {
	Execute(ctx context.Context, p plan.Plan) error
}
