package runtime

import "context"

// Runner executes a launcher script.
type Runner interface {
	Run(ctx context.Context, path string) error
}
