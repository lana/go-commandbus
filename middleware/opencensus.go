package middleware

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	commandbus "github.com/lana/go-commandbus"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
)

var (
	// Measures
	ocExecCount = stats.Int64("executions", "Number of command executions", stats.UnitDimensionless)
	ocErrCount  = stats.Int64("errors", "Number of command errors", stats.UnitDimensionless)

	// Tag keys
	commandName = tag.MustNewKey("command.name")

	// Views
	OpenCensusExecCount = &view.View{
		Name:        "execution_count",
		Measure:     ocExecCount,
		Description: "Number of command executions",
		Aggregation: view.Count(),
	}

	OpenCensusErrCount = &view.View{
		Name:        "error_count",
		Measure:     ocErrCount,
		Description: "Number of command error",
		Aggregation: view.Count(),
	}

	OpenCensusExecCountByCommand = &view.View{
		Name:        "execution_count_by_command",
		Measure:     ocExecCount,
		Description: "Number of command executions by command",
		TagKeys:     []tag.Key{commandName},
		Aggregation: view.Count(),
	}

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
				return errors.Wrap(err, "unable to create the tag")
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
