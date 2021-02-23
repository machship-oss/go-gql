package gql

type QueryContainer struct {
	QueryOps     []*QueryOp
	GetOps       []*GetOp
	AggregateOps []*AggregateOp

	queryPart
}

func (cont *QueryContainer) generateQuery() (err error) {
	// process the operations in the container to get their queries
	for _, qryOp := range cont.QueryOps {
		err = cont.generateSubQuery(qryOp)
		if err != nil {
			return newErr(ErrOperationCouldNotBeGenerated, err)
		}
	}

	for _, getOp := range cont.GetOps {
		err = cont.generateSubQuery(getOp)
		if err != nil {
			return newErr(ErrOperationCouldNotBeGenerated, err)
		}
	}

	for _, aggOp := range cont.AggregateOps {
		err = cont.generateSubQuery(aggOp)
		if err != nil {
			return newErr(ErrOperationCouldNotBeGenerated, err)
		}
	}

	cont.queryPart.buildQuery(false)
	return
}
