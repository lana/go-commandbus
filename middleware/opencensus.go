package middleware

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	commandbus "github.com/lana/go-commandbus"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"

	"go.lana.tech/errors"
)

var (
	// Measures
	ocExecCount = stats.Int64("executions", "Number of command executions", stats.UnitDimensionless)
	ocErrCount  = stats.Int64("errors", "Number of command errors", stats.UnitDimensionless)

	// Tag keys
	commandName = tag.MustNewKey("command.name")

	// OpenCensusExecCount view.
	OpenCensusExecCount = &view.View{
		Name:        "execution_count",
		Measure:     ocExecCount,
		Description: "Number of command executions",
		Aggregation: view.Count(),
	}

	// OpenCensusErrCount view.
	OpenCensusErrCount = &view.View{
		Name:        "error_count",
		Measure:     ocErrCount,
		Description: "Number of command error",
		Aggregation: view.Count(),
	}

	// OpenCensusExecCountByCommand view.
	OpenCensusExecCountByCommand = &view.View{
		Name:        "execution_count_by_command",
		Measure:     ocExecCount,
		Description: "Number of command executions by command",
		TagKeys:     []tag.Key{commandName},
		Aggregation: view.Count(),
	}

	// OpenCensusErrCountByCommand view.
	OpenCensusErrCountByCommand = &view.View{
		Name:        "error_count_by_command",
		Measure:     ocErrCount,
		Description: "Number of command error by command",
		TagKeys:     []tag.Key{commandName},
		Aggregation: view.Count(),
	}
)

// OpenCensusConfig defines the config for OpenCensus middleware.
type OpenCensusConfig struct {
	// Views it is a OpenCensus views list.
	Views []*view.View
}

// DefaultOpenCensusConfig is the default OpenCensus middleware config.
var DefaultOpenCensusConfig = OpenCensusConfig{
	Views: []*view.View{
		OpenCensusExecCount,
		OpenCensusErrCount,
		OpenCensusExecCountByCommand,
		OpenCensusErrCountByCommand,
	},
}

// OpenCensus returns a middleware that collect the command metrics.
func OpenCensus() commandbus.MiddlewareFunc {
	return OpenCensusWithConfig(DefaultOpenCensusConfig)
}

// OpenCensusWithConfig returns a OpenCensus middleware with config.
// See: `OpenCensus()`.
func OpenCensusWithConfig(cfg OpenCensusConfig) commandbus.MiddlewareFunc {
	// Defaults
	if len(cfg.Views) == 0 {
		cfg.Views = DefaultOpenCensusConfig.Views
	}

	if err := view.Register(cfg.Views...); err != nil {
		panic(fmt.Sprintf("middleware error when register views: %v", err))
	}

	return func(next commandbus.HandlerFunc) commandbus.HandlerFunc {
		return func(ctx context.Context, cmd interface{}) error {
			name := reflect.TypeOf(cmd).Elem().Name()

			ctx, span := trace.StartSpan(ctx, strings.ToLower(name))
			defer span.End()

			ctx, err := tag.New(ctx, tag.Upsert(commandName, name))
			if err != nil {
				return errors.Errorf("unable to create the tag: %w", err)
			}

			stats.Record(ctx, ocExecCount.M(1))

			if err := next(ctx, cmd); err != nil {
				span.SetStatus(
					trace.Status{
						Code:    trace.StatusCodeUnknown,
						Message: err.Error(),
					},
				)

				stats.Record(ctx, ocErrCount.M(1))

				return err
			}

			return nil
		}
	}
}
