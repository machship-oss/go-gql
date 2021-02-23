package gql

import (
	"testing"
)

func Test_AddOp(t *testing.T) {
	tests := []struct {
		queryTest
		Name         string
		ResultObject IMutationResult
		Arguments    IIsAdd
		ExpectedErr  error
	}{
		{
			Name:         "Check expected query and variableJSON from valid request with variables",
			ResultObject: &testQueryableMutationResult{},
			Arguments: &AddtestQueryableInputs{
				&AddtestQueryableInput{
					AddBaseInput: AddBaseInput{
						DateCreated: NewTime(testTime),
					},
					testQueryableFields: testQueryableFields{
						ShouldIgnore: NewString("IGNORED"),
						TString:      NewString("something"),
					},
					SomeRef: []*BaseRef{
						{
							ID: NewID("ABCD"),
						},
					},
				},
			},
			ExpectedErr: nil,
			queryTest: queryTest{
				ExpectedQuery:         `add_0:addtestQueryable(input:$input_add_0){numUids,testQueryable{id,dateCreatedUTC,__typename,tString,shouldRecurse{id,dateCreatedUTC,__typename,tString}}}`,
				ExpectedVariablesJSON: `"input_add_0":[{"dateCreated":"2000-01-01T00:00:00Z","tString":"something","someRef":[{"id":"ABCD"}]}]`,
			},
		},
	}

	for _, test := range tests {
		test.OverallTestName = "Test_AddOp"
		op := AddOp{
			ResultObject: test.ResultObject,
			Arguments:    test.Arguments,
		}
		err := op.generateQuery()
		checkErrs(t, test.Name, test.OverallTestName, test.ExpectedErr, err)
		checkQueries(t, test.Name, op.queryPart, test.queryTest)
	}
}
