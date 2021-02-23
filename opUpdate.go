package gql

type UpdateOp struct {
	ResultObject IMutationResult
	Arguments    IIsUpdate

	queryPart
}

func (op *UpdateOp) getQuery() string {
	return op.query
}

func (op *UpdateOp) getQueryParams() []string {
	return op.outerQueryParams
}

func (op *UpdateOp) getQueryVariablesJSON() []string {
	return op.argumentsJSON
}

func (op *UpdateOp) generateQuery(isTestMode ...bool) (err error) {
	qsb := op.prepareQuery("update", op.ResultObject.MutationName(), isTestMode...)

	//deal with any argument inputs
	if op.Arguments == nil {
		return newErr(ErrArgumentsCannotBeEmpty, nil)
	}

	err = op.addParamFieldsAndJson(op.Arguments, at_Update)
	if err != nil {
		return err
	}

	op.finaliseQuery(qsb, op.ResultObject)
	return
}
