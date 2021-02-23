package gql

type AddOp struct {
	ResultObject IMutationResult
	Arguments    IIsAdd

	queryPart
}

func (op *AddOp) getQuery() string {
	return op.query
}

func (op *AddOp) getQueryParams() []string {
	return op.outerQueryParams
}

func (op *AddOp) getQueryVariablesJSON() []string {
	return op.argumentsJSON
}

func (op *AddOp) generateQuery(isTestMode ...bool) (err error) {
	qsb := op.prepareQuery("add", op.ResultObject.MutationName(), isTestMode...)

	//deal with any argument inputs
	if op.Arguments == nil {
		return newErr(ErrArgumentsCannotBeEmpty, nil)
	}

	err = op.addParamFieldsAndJson(op.Arguments, at_Add, op.Arguments.AddName())
	if err != nil {
		return err
	}

	op.finaliseQuery(qsb, op.ResultObject)
	return
}
