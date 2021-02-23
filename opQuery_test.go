package gql

import (
	"testing"
)

func Test_QueryOp(t *testing.T) {
	tests := []struct {
		queryTest
		Name         string
		ResultObject IMultiResult
		Arguments    *QueryArguments
		ExpectedErr  error
	}{
		{
			Name:         "Check expected query from valid request with no variables",
			ResultObject: &testQueryableMultiResult{},
			ExpectedErr:  nil,
			queryTest: queryTest{
				ExpectedQuery:       "query_0:querytestQueryable{id,dateCreatedUTC,__typename,tString,shouldRecurse{id,dateCreatedUTC,__typename,tString}}",
				IgnoreVariablesJSON: true,
			},
		},
		{
			Name:         "Should force JSON error",
			ResultObject: &testQueryableMultiResult{},
			Arguments: &QueryArguments{
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
			ResultObject: &testQueryableMultiResult{},
			Arguments: &QueryArguments{
				Filter: &testQueryableFilter{
					BaseFilter:   BaseFilter{ID: NewID("ABCD")},
					TString:      &StringFilter{Equal: NewString("XYZ")},
					ShouldIgnore: &StringFilter{Equal: NewString("BAD")},
				},
				Order: &testQueryableOrder{
					Asc: OC_SomeFakeField,
					Then: &testQueryableOrder{
						Desc: testQueryableOrderChoice(OC_BaseDateCreated),
					},
				},
				First:  NewInt(10),
				Offset: NewInt(20),
			},
			ExpectedErr: nil,
			queryTest: queryTest{
				ExpectedQuery:         "query_0:querytestQueryable(filter:$filter_query_0,order:$order_query_0,first:$first_query_0,offset:$offset_query_0){id,dateCreatedUTC,__typename,tString,shouldRecurse{id,dateCreatedUTC,__typename,tString}}",
				ExpectedVariablesJSON: `"filter_query_0":{"id":"ABCD","tString":{"eq":"XYZ"}},"order_query_0":{"asc":"someFakeField","then":{"desc":"dateCreated"}},"first_query_0":10,"offset_query_0":20`,
			},
		},
		{
			Name:         "Check expected query and variableJSON from valid request with variables with relationships",
			ResultObject: &testQueryableMultiResult{},
			Arguments: &QueryArguments{Filter: &testQueryableFilter{
				BaseFilter: BaseFilter{ID: NewID("ABCD")},
				TString:    &StringFilter{Equal: NewString("XYZ")},
			}},
			ExpectedErr: nil,
			queryTest: queryTest{
				ExpectedQuery:         "query_0:querytestQueryable(filter:$filter_query_0){id,dateCreatedUTC,__typename,tString,shouldRecurse{id,dateCreatedUTC,__typename,tString}}",
				ExpectedVariablesJSON: `"filter_query_0":{"id":"ABCD","tString":{"eq":"XYZ"}}`,
			},
		},
	}

	for _, test := range tests {
		test.OverallTestName = "Test_QueryOp"
		op := QueryOp{
			ResultObject: test.ResultObject,
			Arguments:    test.Arguments,
		}
		err := op.generateQuery()
		checkErrs(t, test.Name, test.OverallTestName, test.ExpectedErr, err)
		checkQueries(t, test.Name, op.queryPart, test.queryTest)
	}
}
