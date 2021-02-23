package gql

import "testing"

func Test_DeleteOp(t *testing.T) {
	tests := []struct {
		queryTest
		Name         string
		ResultObject IMutationResult
		Arguments    IIsDelete
		ExpectedErr  error
	}{
		{
			Name:         "Check expected query and variableJSON from valid request with variables",
			ResultObject: &testQueryableDeleteResult{},
			Arguments: &testQueryableDeleteArguments{
				Filter: &testQueryableFilter{
					BaseFilter:   BaseFilter{ID: NewID("ABCD")},
					TString:      &StringFilter{Equal: NewString("XYZ")},
					ShouldIgnore: &StringFilter{Equal: NewString("BAD")},
				},
			},
			ExpectedErr: nil,
			queryTest: queryTest{
				ExpectedQuery:         `delete_0:deletetestQueryable(filter:$filter_delete_0){numUids,msg,testQueryable{id,dateCreatedUTC,__typename,tString,shouldRecurse{id,dateCreatedUTC,__typename,tString}}}`,
				ExpectedVariablesJSON: `"filter_delete_0":{"id":"ABCD","tString":{"eq":"XYZ"}}`,
			},
		},
	}

	for _, test := range tests {
		test.OverallTestName = "Test_UpdateOp"
		op := DeleteOp{
			ResultObject: test.ResultObject,
			Arguments:    test.Arguments,
		}
		err := op.generateQuery()
		checkErrs(t, test.Name, test.OverallTestName, test.ExpectedErr, err)
		checkQueries(t, test.Name, op.queryPart, test.queryTest)
	}
}
