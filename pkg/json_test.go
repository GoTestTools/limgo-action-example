package pkg

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestToJSON_Succeeds(t *testing.T) {
	cfg := CoverageConfig{
		CoverageThreshold: CoverageThreshold{
			Global: &Threshold{
				Statements: 100,
				Functions:  100,
			},
			Matcher: &Matcher{
				"pkg/coverage/line.go": Threshold{
					Statements: 10,
					Functions:  20,
				},
			},
		},
		Excludes: []string{"vendor/.*"},
	}

	cfgBuffer := &bytes.Buffer{}
	err := cfg.ToJSON(cfgBuffer)
	if err != nil {
		t.Fatalf("Unexpected error writing config struct to string: %v", err)
	}

	expectedRawCfg := `{
	"coverageThreshold": {
		"global": {
			"statements": 100,
			"functions": 100
		},
		"matcher": {
			"pkg/coverage/line.go": {
				"statements": 10,
				"functions": 20
			}
		}
	},
	"excludes": [
		"vendor/.*"
	]
}
`

	if diff := cmp.Diff(expectedRawCfg, cfgBuffer.String()); diff != "" {
		t.Fatalf("Detected difference in parsed config: %s", diff)
	}
}

func TestFromJSONString_Succeeds(t *testing.T) {
	rawCfg := `
	{
		"coverageThreshold": {
			"global": {
				"statements": 100,
				"functions": 100
			},
			"matcher": {
				"pkg/coverage/line.go": {
					"statements": 10,
					"functions": 20
				}
			}
		},
		"excludes": [
			"vendor/.*"
		]
	}
	`

	parsedCfg, err := FromJSONString(strings.NewReader(rawCfg))
	if err != nil {
		t.Fatalf("Unexpected error parsing config string to struct: %v", err)
	}

	expectedCfg := CoverageConfig{
		CoverageThreshold: CoverageThreshold{
			Global: &Threshold{
				Statements: 100,
				Functions:  100,
			},
			Matcher: &Matcher{
				"pkg/coverage/line.go": Threshold{
					Statements: 10,
					Functions:  20,
				},
			},
		},
		Excludes: []string{"vendor/.*"},
	}

	if diff := cmp.Diff(expectedCfg, parsedCfg); diff != "" {
		t.Fatalf("Detected difference in parsed config: %s", diff)
	}
}

func TestFromJSONString_FailsDueToInvalidJSON(t *testing.T) {
	cfgMissingBracket := `
	{
		"coverageThreshold": {
			"global": {
				"statements": 100,
				"functions": 100
			},
		}
	`

	parsedCfg, err := FromJSONString(strings.NewReader(cfgMissingBracket))
	if err == nil {
		t.Fatalf("Expected error when parsing invalid config string, but got cfg: %v", parsedCfg)
	}
}
