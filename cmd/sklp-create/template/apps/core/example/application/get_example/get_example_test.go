package get_example

import (
	"testing"

	"github.com/cucumber/godog"
)

// TestFeatures runs the Gherkin scenarios in get_example.feature against the
// RegisterSteps definitions.
func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: RegisterSteps,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"get_example.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
