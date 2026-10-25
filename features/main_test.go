package features

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/cucumber/godog"
)

// doneTags lists the tags of merged tasks; each task extends it with its own tag.
// GODOG_TAGS overrides it when set.
//
// Convention for step definitions: one initXxx(sc *godog.ScenarioContext) per
// feature file, in its own *_steps_test.go file and called from initializeScenario.
// Step patterns are anchored (^...$), and a step text already defined by another
// file is reused, never registered twice.
var doneTags = []string{"@T-01-02", "@T-01-03", "@T-01-04", "@T-01-05"}

func tags() string {
	if t := os.Getenv("GODOG_TAGS"); t != "" {
		return t
	}
	if len(doneTags) == 0 {
		return "@none" // no merged task yet: select no scenario
	}
	return strings.Join(doneTags, ",")
}

func initializeScenario(sc *godog.ScenarioContext) {
	initInput(sc)
	initTiming(sc)
	initLetterbox(sc)
	initDisplay(sc)
}

// TestFeatures runs the Gherkin scenarios of the merged tasks.
func TestFeatures(t *testing.T) {
	scenarios := 0
	suite := godog.TestSuite{
		ScenarioInitializer: func(sc *godog.ScenarioContext) {
			sc.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
				scenarios++
				return ctx, nil
			})
			initializeScenario(sc)
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"."},
			Strict:   true,
			TestingT: t,
			Tags:     tags(),
		},
	}
	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
	if tags() != "@none" && scenarios == 0 {
		t.Fatalf("tag selection %q ran zero scenarios", tags())
	}
}
