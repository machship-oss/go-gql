package gql

type GetOp struct {
	ResultObject ISingleResult
	Arguments    IIsGet

	queryPart
}

func (op *GetOp) getQuery() string {
	return op.query
}

func (op *GetOp) getQueryParams() []string {
	return op.outerQueryParams
}

func (op *GetOp) getQueryVariablesJSON() []string {
	return op.argumentsJSON
}

func (op *GetOp) generateQuery(isTestMode ...bool) (err error) {
	qsb := op.prepareQuery("get", op.ResultObject.SingleResultName(), isTestMode...)

	//deal with any argument inputs
	if op.Arguments == nil {
		return newErr(ErrArgumentsCannotBeEmpty, nil)
	}

	err = op.addParamFieldsAndJson(op.Arguments, at_Get)
	if err != nil {
		return err
	}

	op.finaliseQuery(qsb, op.ResultObject)
	return
}
