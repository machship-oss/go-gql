package gql

type QueryOp struct {
	ResultObject IMultiResult
	Arguments    *QueryArguments

	queryPart
}

type QueryArguments struct {
	Filter IFilter
	Order  IOrder
	First  *Int
	Offset *Int
}

func (op *QueryOp) getQuery() string {
	return op.query
}

func (op *QueryOp) getQueryParams() []string {
	return op.outerQueryParams
}

func (op *QueryOp) getQueryVariablesJSON() []string {
	return op.argumentsJSON
}

func (op *QueryOp) generateQuery(isTestMode ...bool) (err error) {
	qsb := op.prepareQuery("query", op.ResultObject.MultiResultName(), isTestMode...)

	//deal with any argument inputs
	if op.Arguments != nil {
		if fltr := op.Arguments.Filter; fltr != nil {
			err = op.addParamsAndJson(fltr, fltr.GetName(), "filter")
			if err != nil {
				return err
			}
		}

		if ordr := op.Arguments.Order; ordr != nil {
			err = op.addParamsAndJson(ordr, ordr.GetName(), "order")
			if err != nil {
				return err
			}
		}

		if frst := op.Arguments.First; frst != nil {
			err = op.addParamsAndJson(frst, frst.GetName(), "first")
			if err != nil {
				return err
			}
		}

		if ofst := op.Arguments.Offset; ofst != nil {
			err = op.addParamsAndJson(ofst, ofst.GetName(), "offset")
			if err != nil {
				return err
			}
		}
	}

	op.finaliseQuery(qsb, op.ResultObject)
	return
}
