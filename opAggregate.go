package gql

type AggregateOp struct {
	ResultObject IAggregateResult
	Arguments    *AggregateArguments

	queryPart
}

type AggregateArguments struct {
	Filter IFilter
}

func (op *AggregateOp) getQuery() string {
	return op.query
}

func (op *AggregateOp) getQueryParams() []string {
	return op.outerQueryParams
}

func (op *AggregateOp) getQueryVariablesJSON() []string {
	return op.argumentsJSON
}

func (op *AggregateOp) generateQuery(isTestMode ...bool) (err error) {
	qsb := op.prepareQuery("aggregate", op.ResultObject.AggregateResultName(), isTestMode...)

	//deal with any argument inputs
	if op.Arguments != nil {
		if fltr := op.Arguments.Filter; fltr != nil {
			err = op.addParamsAndJson(fltr, fltr.GetName(), "filter")
			if err != nil {
				return err
			}
		}
	}

	op.finaliseQuery(qsb, op.ResultObject)
	return
}
