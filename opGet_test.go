package gql

import (
	"testing"
)

func Test_GetOp(t *testing.T) {
	tests := []struct {
		queryTest

		Name         string
		ResultObject ISingleResult
		Arguments    IIsGet
		ExpectedErr  error
	}{
		{
			Name:         "Check expected query and variableJSON from valid request with variables",
			ResultObject: &testQueryable{},
			Arguments: &testQueryableGet{
				ID:      NewID("ABCD"),
				TString: NewString("XYZ"),
			},
			ExpectedErr: nil,
			queryTest: queryTest{
				ExpectedQuery:         `get_0:gettestQueryable(id:$id_get_0,tString:$tString_get_0){id,dateCreatedUTC,__typename,tString,shouldRecurse{id,dateCreatedUTC,__typename,tString}}`,
				ExpectedVariablesJSON: `"id_get_0":"ABCD","tString_get_0":"XYZ"`,
			},
		},
	}

	for _, test := range tests {
		test.OverallTestName = "Test_GetOp"
		op := GetOp{
			ResultObject: test.ResultObject,
			Arguments:    test.Arguments,
		}
		err := op.generateQuery()
		checkErrs(t, test.Name, test.OverallTestName, test.ExpectedErr, err)
		checkQueries(t, test.Name, op.queryPart, test.queryTest)
	}
}
