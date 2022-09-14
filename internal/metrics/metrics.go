package metrics

import (
	"context"
	"fmt"
	"os"

	"github.com/brexhq/substation/internal/errors"
)

// unsupportedDestination is returned when an unsupported Metrics destination is referenced in Factory.
const unsupportedDestination = errors.Error("unsupportedDestination")

// referenced across all metrics generators
var metricsDestination string
var metricsApplication string

// used when generating metrics from AWS Lambda functions
var metricsAWSLambdaFunctionName string

func init() {
	// determines if metrics should be generated across the application. the value from this environment variable is used to retrieve the metrics destination from the Factory.
	if m, ok := os.LookupEnv("SUBSTATION_METRICS"); ok {
		metricsDestination = m
	}

	metricsApplication = "Substation"

	metricsAWSLambdaFunctionName, _ = os.LookupEnv("AWS_LAMBDA_FUNCTION_NAME")
}

// Data contains a metric that can be sent to external services.
type Data struct {
	// Contextual information related to the metric. If the external service accepts key-value pairs (e.g., identifiers, tags), then this is passed directly to the service.
	Attributes map[string]string

	// A short name that describes the metric. This is passed directly to the external service and should use the upper camel case (UpperCamelCase) naming convention.
	Name string

	// The metric data point. This value is converted to the correct data type before being sent to the external service.
	Value interface{}
}

// AddAttributes is a convenience method for adding attributes to a metric.
func (d *Data) AddAttributes(attr map[string]string) {
	if d.Attributes == nil {
		d.Attributes = make(map[string]string)
	}

	for key, val := range attr {
		d.Attributes[key] = val
	}
}

// Generator is an interface for creating a metric and sending it to an external service.
type Generator interface {
	Generate(context.Context, Data) error
}

// Generate is a convenience function that encapsulates the Factory and creates a metric. If the SUBSTATION_METRICS environment variable is not set, then no metrics are created.
func Generate(ctx context.Context, data Data) error {
	if metricsDestination == "" {
		return nil
	}

	gen, err := Factory(metricsDestination)
	if err != nil {
		return err
	}

	if err := gen.Generate(ctx, data); err != nil {
		return err
	}

	return nil
}

// Factory returns a configured metrics Generator. This is the recommended method for retrieving ready-to-use Generators.
func Factory(destination string) (Generator, error) {
	switch destination {
	case "AWS_CLOUDWATCH_EMBEDDED_METRICS":
		var m AWSCloudWatchEmbeddedMetrics
		return m, nil
	default:
		return nil, fmt.Errorf("metrics destination %s: %v", destination, unsupportedDestination)
	}
}