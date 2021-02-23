package gql

import "testing"

func Test_UpdateOp(t *testing.T) {
	tests := []struct {
		queryTest
		Name         string
		ResultObject IMutationResult
		Arguments    IIsUpdate
		ExpectedErr  error
	}{
		{
			Name:         "Check expected query and variableJSON from valid request with variables",
			ResultObject: &testQueryableMutationResult{},
			Arguments: &updatetestQueryableInput{
				Filter: &testQueryableFilter{
					BaseFilter:   BaseFilter{ID: NewID("ABCD")},
					TString:      &StringFilter{Equal: NewString("XYZ")},
					ShouldIgnore: &StringFilter{Equal: NewString("BAD")},
				},
				Set: &testQueryablePatch{
					BasePatch: BasePatch{
						DateCreated: NewTime(testTime),
					},
				},
			},
			ExpectedErr: nil,
			queryTest: queryTest{
				ExpectedQuery:         `update_0:updatetestQueryable(input:$input_update_0){numUids,testQueryable{id,dateCreatedUTC,__typename,tString,shouldRecurse{id,dateCreatedUTC,__typename,tString}}}`,
				ExpectedVariablesJSON: `"input_update_0":{"filter":{"id":"ABCD","tString":{"eq":"XYZ"}},"set":{"dateCreated":"2000-01-01T00:00:00Z"}}`,
			},
		},
	}

	for _, test := range tests {
		test.OverallTestName = "Test_UpdateOp"
		op := UpdateOp{
			ResultObject: test.ResultObject,
			Arguments:    test.Arguments,
		}
		err := op.generateQuery()
		checkErrs(t, test.Name, test.OverallTestName, test.ExpectedErr, err)
		checkQueries(t, test.Name, op.queryPart, test.queryTest)
	}
}
