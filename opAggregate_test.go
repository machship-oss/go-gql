package gql

import (
	"testing"
)

func Test_AggregateOp(t *testing.T) {
	tests := []struct {
		queryTest
		Name         string
		ResultObject IAggregateResult
		Arguments    *AggregateArguments
		ExpectedErr  error
	}{
		{
			Name:         "Check expected query from valid request with no variables",
			ResultObject: &testQueryableAggregateResult{},
			ExpectedErr:  nil,
			queryTest: queryTest{
				ExpectedQuery:       "aggregate_0:aggregatetestQueryable{count}",
				IgnoreVariablesJSON: true,
			},
		},
		{
			Name:         "Should force JSON error",
			ResultObject: &testQueryableAggregateResult{},
			Arguments: &AggregateArguments{
				Filter: &shouldForceJSONError{
					FailMe: make(chan int),
				},
			},
			ExpectedErr: ErrArgumentCouldNotBeEncodedAsJSON,
			queryTest: queryTest{
				IgnoreQuery:         true,
				IgnoreVariablesJSON: true,
			},
		},
		{
			Name:         "Check expected query and variableJSON from valid request with variables",
			ResultObject: &testQueryableAggregateResult{},
			Arguments: &AggregateArguments{
				Filter: &testQueryableFilter{
					BaseFilter:   BaseFilter{ID: NewID("ABCD")},
					TString:      &StringFilter{Equal: NewString("XYZ")},
					ShouldIgnore: &StringFilter{Equal: NewString("BAD")},
				},
			},
			ExpectedErr: nil,
			queryTest: queryTest{
				ExpectedQuery:         "aggregate_0:aggregatetestQueryable(filter:$filter_aggregate_0){count}",
				ExpectedVariablesJSON: `"filter_aggregate_0":{"id":"ABCD","tString":{"eq":"XYZ"}}`,
			},
		},
		{
			//todo: fix this so that it actually uses relationships
			Name:         "Check expected query and variableJSON from valid request with variables with relationships",
			ResultObject: &testQueryableAggregateResult{},
			Arguments: &AggregateArguments{Filter: &testQueryableFilter{
				BaseFilter: BaseFilter{ID: NewID("ABCD")},
				TString:    &StringFilter{Equal: NewString("XYZ")},
			}},
			ExpectedErr: nil,
			queryTest: queryTest{
				ExpectedQuery:         "aggregate_0:aggregatetestQueryable(filter:$filter_aggregate_0){count}",
				ExpectedVariablesJSON: `"filter_aggregate_0":{"id":"ABCD","tString":{"eq":"XYZ"}}`,
			},
		},
		{
			Name:         "Check expected query and variableJSON from valid request with variables with has filter",
			ResultObject: &testQueryableAggregateResult{},
			Arguments: &AggregateArguments{Filter: &testQueryableFilter{
				BaseFilter: BaseFilter{ID: NewID("ABCD")},
				Has:        HC_testQueryable_DateCreatedUTC,
			}},
			ExpectedErr: nil,
			queryTest: queryTest{
				ExpectedQuery:         "aggregate_0:aggregatetestQueryable(filter:$filter_aggregate_0){count}",
				ExpectedVariablesJSON: `"filter_aggregate_0":{"id":"ABCD","has":"dateCreatedUTC"}`,
			},
		},
	}

	for _, test := range tests {
		test.OverallTestName = "Test_AggregateOp"
		op := AggregateOp{
			ResultObject: test.ResultObject,
			Arguments:    test.Arguments,
		}
		err := op.generateQuery()
		checkErrs(t, test.Name, test.OverallTestName, test.ExpectedErr, err)
		checkQueries(t, test.Name, op.queryPart, test.queryTest)
	}
}
