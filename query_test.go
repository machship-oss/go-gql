package gql

import (
	"testing"
)

func Test_QueryColl(t *testing.T) {
	tests := []struct {
		queryTest
		Name         string
		QueryOps     []*QueryOp
		GetOps       []*GetOp
		AggregateOps []*AggregateOp
		ExpectedErr  error
	}{
		{
			Name: "Check expected query from valid request with no variables",
			QueryOps: []*QueryOp{{
				ResultObject: &testQueryableMultiResult{},
			}},
			ExpectedErr: nil,
			queryTest: queryTest{
				ExpectedQuery:        "query{query_0:querytestQueryable{id,dateCreatedUTC,__typename,tString,shouldRecurse{id,dateCreatedUTC,__typename,tString}}}",
				ExpectedArgumentJSON: "{}",
				IgnoreVariablesJSON:  true,
			},
		},
		{
			Name: "Check expected query and variableJSON from valid QueryOps request with variables",
			QueryOps: []*QueryOp{
				{
					ResultObject: &testQueryableMultiResult{},
					Arguments: &QueryArguments{
						Filter: &testQueryableFilter{
							BaseFilter: BaseFilter{ID: NewID("ABCD")},
						},
					},
				},
			},
			ExpectedErr: nil,
			queryTest: queryTest{
				ExpectedQuery:        "query($filter_query_0:testQueryableFilter){query_0:querytestQueryable(filter:$filter_query_0){id,dateCreatedUTC,__typename,tString,shouldRecurse{id,dateCreatedUTC,__typename,tString}}}",
				ExpectedArgumentJSON: `{"filter_query_0":{"id":"ABCD"}}`,
				IgnoreVariablesJSON:  true,
			},
		},
		{
			Name: "Check expected query and variableJSON from valid GetOps request with variables",
			GetOps: []*GetOp{
				{
					ResultObject: &testQueryable{},
					Arguments: &testQueryableGet{
						ID:      NewID("ABCD"),
						TString: NewString("XYZ"),
					},
				},
			},
			ExpectedErr: nil,
			queryTest: queryTest{
				ExpectedQuery:        `query($id_get_0:ID,$tString_get_0:String){get_0:gettestQueryable(id:$id_get_0,tString:$tString_get_0){id,dateCreatedUTC,__typename,tString,shouldRecurse{id,dateCreatedUTC,__typename,tString}}}`,
				ExpectedArgumentJSON: `{"id_get_0":"ABCD","tString_get_0":"XYZ"}`,
				IgnoreVariablesJSON:  true,
			},
		},
		{
			Name: "Check expected query and variableJSON from valid AggregateOps request with variables",
			AggregateOps: []*AggregateOp{
				{
					ResultObject: &testQueryableAggregateResult{},
					Arguments: &AggregateArguments{
						Filter: &testQueryableFilter{
							BaseFilter: BaseFilter{ID: NewID("ABCD")},
						},
					},
				},
			},
			ExpectedErr: nil,
			queryTest: queryTest{
				ExpectedQuery:        "query($filter_aggregate_0:testQueryableFilter){aggregate_0:aggregatetestQueryable(filter:$filter_aggregate_0){count}}",
				ExpectedArgumentJSON: `{"filter_aggregate_0":{"id":"ABCD"}}`,
				IgnoreVariablesJSON:  true,
			},
		},

		{
			Name: "Check expected query from valid request with no variables",
			QueryOps: []*QueryOp{
				{
					ResultObject: &testQueryableMultiResult{},
					Arguments: &QueryArguments{
						Filter: &testQueryableFilter{
							BaseFilter: BaseFilter{ID: NewID("ABCD")},
						},
					},
				},
			},
			GetOps: []*GetOp{
				{
					ResultObject: &testQueryable{},
					Arguments: &testQueryableGet{
						ID:      NewID("ABCD"),
						TString: NewString("XYZ"),
					},
				},
			},
			ExpectedErr: nil,
			queryTest: queryTest{
				ExpectedQuery:        "query($filter_query_0:testQueryableFilter,$id_get_0:ID,$tString_get_0:String){query_0:querytestQueryable(filter:$filter_query_0){id,dateCreatedUTC,__typename,tString,shouldRecurse{id,dateCreatedUTC,__typename,tString}},get_0:gettestQueryable(id:$id_get_0,tString:$tString_get_0){id,dateCreatedUTC,__typename,tString,shouldRecurse{id,dateCreatedUTC,__typename,tString}}}",
				ExpectedArgumentJSON: `{"filter_query_0":{"id":"ABCD"},"id_get_0":"ABCD","tString_get_0":"XYZ"}`,
				IgnoreVariablesJSON:  true,
			},
		},
	}

	for _, test := range tests {
		test.OverallTestName = "Test_QueryColl"
		op := QueryContainer{
			QueryOps:     test.QueryOps,
			GetOps:       test.GetOps,
			AggregateOps: test.AggregateOps,
		}

		err := op.generateQuery()
		checkErrs(t, test.Name, test.OverallTestName, test.ExpectedErr, err)
		checkQueries(t, test.Name, op.queryPart, test.queryTest)
		checkQueriesColl(t, test.Name, op.queryPart, test.queryTest)
	}
}

func Test_QueryPartIdentifiers(t *testing.T) {
	tests := []struct {
		testBase
		Name                string
		QueryOps            []*QueryOp
		ExpectedIdentifiers []string
		ExpectedErr         error
	}{
		{
			Name: "Check that identifiers count up",
			QueryOps: []*QueryOp{
				{
					ResultObject: &testQueryableMultiResult{},
				},
				{
					ResultObject: &testQueryableMultiResult{},
				},
			},
			ExpectedIdentifiers: []string{"query_0", "query_1"},
			ExpectedErr:         nil,
		},
	}

	for _, test := range tests {
		test.OverallTestName = "Test_QueryPart"

		if len(test.ExpectedIdentifiers) != len(test.QueryOps) {
			t.Errorf("\n[%s] Failed Identifier Check '%s'; number of QueryOps (%d) doesn't match number of ExpectedIdentifiers (%d)", test.OverallTestName, test.Name, len(test.QueryOps), len(test.ExpectedIdentifiers))
			t.FailNow()
		}

		for idx, op := range test.QueryOps {
			err := op.generateQuery(false)
			checkErrs(t, test.Name, test.OverallTestName, test.ExpectedErr, err)
			if test.ExpectedIdentifiers[idx] != op.identifier {
				t.Errorf("\n[%s] Failed Identifier Check '%s'; expected: '%s', got: '%s'", test.OverallTestName, test.Name, test.ExpectedIdentifiers[idx], op.identifier)
			}
		}
	}
}
