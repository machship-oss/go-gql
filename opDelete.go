package gql

type DeleteOp struct {
	ResultObject IMutationResult
	Arguments    IIsDelete

	queryPart
}

// type DeleteArguments struct {
// 	Filter IFilter `json:"filter"`
// }

// func (t *DeleteArguments) IsArg() {}

func (op *DeleteOp) getQuery() string {
	return op.query
}

func (op *DeleteOp) getQueryParams() []string {
	return op.outerQueryParams
}

func (op *DeleteOp) getQueryVariablesJSON() []string {
	return op.argumentsJSON
}

func (op *DeleteOp) generateQuery(isTestMode ...bool) (err error) {
	qsb := op.prepareQuery("delete", op.ResultObject.MutationName(), isTestMode...)

	//deal with any argument inputs
	// if op.Arguments != nil {
	// 	if fltr := op.Arguments.Filter; fltr != nil {
	// 		err = op.addParamsAndJson(fltr, fltr.Name(), "filter")
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	err = op.addParamFieldsAndJson(op.Arguments.DeleteFilter(), at_Delete, op.Arguments.DeleteName())
	if err != nil {
		return err
	}

	op.finaliseQuery(qsb, op.ResultObject)
	return
}
