package gql

import (
	"testing"
)

func Test_MutationColl(t *testing.T) {
	tests := []struct {
		queryTest
		Name        string
		AddOps      []*AddOp
		ExpectedErr error
	}{
		{
			Name: "Check expected query and variableJSON from valid request with variables",
			AddOps: []*AddOp{
				{
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
				},
			},
			ExpectedErr: nil,
			queryTest: queryTest{
				//todo: check how to know when to add the exclamation marks
				ExpectedQuery:        `mutation($input_add_0:[AddtestQueryableInput!]!){add_0:addtestQueryable(input:$input_add_0){numUids,testQueryable{id,dateCreatedUTC,__typename,tString,shouldRecurse{id,dateCreatedUTC,__typename,tString}}}}`,
				ExpectedArgumentJSON: `{"input_add_0":[{"dateCreated":"2000-01-01T00:00:00Z","tString":"something","someRef":[{"id":"ABCD"}]}]}`,
				IgnoreVariablesJSON:  true,
			},
		},
	}

	for _, test := range tests {
		test.OverallTestName = "Test_MutationColl"
		op := MutationContainer{
			AddOps: test.AddOps,
		}

		err := op.generateQuery()
		checkErrs(t, test.Name, test.OverallTestName, test.ExpectedErr, err)
		checkQueries(t, test.Name, op.queryPart, test.queryTest)
		checkQueriesColl(t, test.Name, op.queryPart, test.queryTest)
	}
}
