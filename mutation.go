package gql

type MutationContainer struct {
	AddOps    []*AddOp
	UpdateOps []*UpdateOp
	DeleteOps []*DeleteOp

	queryPart
}

func (cont *MutationContainer) generateQuery() (err error) {
	for _, addOp := range cont.AddOps {
		err = cont.generateSubQuery(addOp)
		if err != nil {
			return newErr(ErrOperationCouldNotBeGenerated, err)
		}
	}

	for _, updateOp := range cont.UpdateOps {
		err = cont.generateSubQuery(updateOp)
		if err != nil {
			return newErr(ErrOperationCouldNotBeGenerated, err)
		}
	}

	for _, deleteOp := range cont.DeleteOps {
		err = cont.generateSubQuery(deleteOp)
		if err != nil {
			return newErr(ErrOperationCouldNotBeGenerated, err)
		}
	}

	cont.queryPart.buildQuery(true)
	return
}
