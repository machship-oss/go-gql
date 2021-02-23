package gql

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	testMode = true
	code := m.Run()
	os.Exit(code)
}

var (
	testTime = time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC)
)

func checkErrs(t *testing.T, overallTestName, testName string, expected, actual error) {
	if actual != nil && expected == nil {
		t.Errorf("\n[%s] Failed Error Check '%s'; expected: '%v', got: '%v'", overallTestName, testName, expected, actual)
	} else if actual == nil && expected != nil {
		t.Errorf("\n[%s] Failed Error Check '%s'; expected: '%v', got: '%v'", overallTestName, testName, expected, actual)
	} else if actual != nil && expected.Error() != actual.Error() {
		t.Errorf("\n[%s] Failed Error Check '%s'; expected: '%v', got: '%v'", overallTestName, testName, expected, actual)
	}
}

type testBase struct {
	OverallTestName string
}

type queryTest struct {
	testBase
	ExpectedQuery         string
	IgnoreQuery           bool
	ExpectedVariablesJSON string
	IgnoreVariablesJSON   bool
	ExpectedArgumentJSON  string
	IgnoreArgumentJSON    bool
}

func checkQueries(t *testing.T, testName string, qryPart queryPart, test queryTest) {
	if !test.IgnoreQuery && qryPart.query != test.ExpectedQuery {
		t.Errorf("\n[%s] Failed Query Check '%s';\n\nexpected: '%v'\n\ngot: '%v'", test.OverallTestName, testName, test.ExpectedQuery, qryPart.query)
	}

	if got := strings.Join(qryPart.argumentsJSON, ","); !test.IgnoreVariablesJSON && got != test.ExpectedVariablesJSON {
		t.Errorf("\n[%s] Failed Variables JSON Check '%s';\n\nexpected: '%v'\n\ngot: '%v'", test.OverallTestName, testName, test.ExpectedVariablesJSON, got)
	}
}

func checkQueriesColl(t *testing.T, testName string, qryPart queryPart, test queryTest) {
	if !test.IgnoreArgumentJSON && qryPart.argumentJSON != test.ExpectedArgumentJSON {
		t.Errorf("\n[%s] Failed Arguments JSON Check '%s';\n\nexpected: '%v'\n\ngot: '%v'", test.OverallTestName, testName, test.ExpectedArgumentJSON, qryPart.argumentJSON)
	}
}
