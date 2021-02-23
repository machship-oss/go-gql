package main

import (
	"github.com/machship-oss/go-gql"
)

type IBase struct {
	*IBaseFields

	*IBaseRelationships
}

func (ar *IBase) SingleResultName() string {
	return "IBase"
}

type IBaseMultiResult []*IBase

func (ar *IBaseMultiResult) MultiResultName() string {
	return "IBase"
}

type IBaseFields struct {
	ID             *gql.ID   `json:"id,omitempty"`
	DateCreatedUTC *gql.Time `json:"dateCreatedUTC,omitempty"`
}

type IBaseGetArgs struct {
	ID             *gql.ID   `json:"id,omitempty"`
	DateCreatedUTC *gql.Time `json:"dateCreatedUTC,omitempty"`
}

func (t *IBaseGetArgs) IsArg() {}

func (t *IBaseGetArgs) IsGet() {}

type IBaseRelationships struct {
}

type UpdateIBaseInput struct {
	Filter *IBaseFilter `json:"filter"`
	Set    *IBasePatch  `json:"set,omitempty"`
	Remove *IBasePatch  `json:"remove,omitempty"`
}

func (t *UpdateIBaseInput) IsArg() {}

func (t *UpdateIBaseInput) IsUpdate() {}

type IBasePatch struct {
	IBaseFields
}

type AddIBaseInput struct {
	IBaseFields
}

type AddIBaseInputs []*AddIBaseInput

func (t *AddIBaseInputs) IsArg() {}

func (t *AddIBaseInputs) AddName() string {
	return "AddIBaseInput"
}

type IBaseAggregateResult struct {
	gql.BaseAggregate
}

func (ar *IBaseAggregateResult) AggregateResultName() string {
	return "IBase"
}

type IBaseMutationResult struct {
	IBases         []*IBase `isResult:"true" json:"iBase,omitempty"`
	NumberAffected gql.Int  `json:"numUids,omitempty"`
}

type IBaseDeleteResult struct {
	IBaseMutationResult
	Message gql.String `json:"msg,omitempty"`
}

func (r IBaseMutationResult) MutationName() string {
	return "IBase"
}

type IBaseDeleteArguments struct {
	Filter *IBaseFilter `json:"filter"`
}

func (t *IBaseDeleteArguments) IsArg() {}

func (t *IBaseDeleteArguments) DeleteName() string {
	return "IBaseFilter"
}

func (t *IBaseDeleteArguments) DeleteFilter() gql.IFilter {
	return t.Filter
}

type IBaseFilter struct {
	ID             *gql.ID   `json:"id,omitempty"`
	DateCreatedUTC *gql.Time `json:"dateCreatedUTC,omitempty"`

	Has IBaseHasChoice `json:"has,omitempty"`
	And []*IBaseFilter `json:"and,omitempty"`
	Or  []*IBaseFilter `json:"or,omitempty"`
	Not []*IBaseFilter `json:"not,omitempty"`
}

func (t *IBaseFilter) GetName() string {
	return "IBaseFilter"
}

func (t *IBaseFilter) IsFilter() {}

func (t *IBaseFilter) IsArg() {}

type IBaseOrderChoice string

const (
	OC_IBase_DateCreatedUTC IBaseOrderChoice = "dateCreatedUTC"
)

type IBaseOrder struct {
	Asc  IBaseOrderChoice `json:"asc,omitempty"`
	Desc IBaseOrderChoice `json:"desc,omitempty"`
	Then *IBaseOrder      `json:"then,omitempty"`
}

func (t *IBaseOrder) GetName() string {
	return "IBaseOrder"
}

func (t *IBaseOrder) IsOrder() {}

type IBaseHasChoice string

const (
	HC_IBase_ID             IBaseHasChoice = "id"
	HC_IBase_DateCreatedUTC IBaseHasChoice = "dateCreatedUTC"
)
