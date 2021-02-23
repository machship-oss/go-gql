package main

import (
	"github.com/machship-oss/go-gql"
)

type ITestInterface struct {
	*ITestInterfaceFields

	*ITestInterfaceRelationships
}

func (ar *ITestInterface) SingleResultName() string {
	return "ITestInterface"
}

type ITestInterfaceMultiResult []*ITestInterface

func (ar *ITestInterfaceMultiResult) MultiResultName() string {
	return "ITestInterface"
}

type ITestInterfaceFields struct {
	ID   *gql.ID     `json:"id,omitempty"`
	Name *gql.String `json:"name,omitempty"`
}

type ITestInterfaceGetArgs struct {
	ID   *gql.ID     `json:"id,omitempty"`
	Name *gql.String `json:"name,omitempty"`
}

func (t *ITestInterfaceGetArgs) IsArg() {}

func (t *ITestInterfaceGetArgs) IsGet() {}

type ITestInterfaceRelationships struct {
	Notes []*Note `hasInverse:"belongsTo" json:"notes,omitempty"`
}

type UpdateITestInterfaceInput struct {
	Filter *ITestInterfaceFilter `json:"filter"`
	Set    *ITestInterfacePatch  `json:"set,omitempty"`
	Remove *ITestInterfacePatch  `json:"remove,omitempty"`
}

func (t *UpdateITestInterfaceInput) IsArg() {}

func (t *UpdateITestInterfaceInput) IsUpdate() {}

type ITestInterfacePatch struct {
	ITestInterfaceFields
	Notes []*gql.BaseRef `json:"notes,omitempty"`
}

type AddITestInterfaceInput struct {
	ITestInterfaceFields
	Notes []*gql.BaseRef `json:"notes,omitempty"`
}

type AddITestInterfaceInputs []*AddITestInterfaceInput

func (t *AddITestInterfaceInputs) IsArg() {}

func (t *AddITestInterfaceInputs) AddName() string {
	return "AddITestInterfaceInput"
}

type ITestInterfaceAggregateResult struct {
	gql.BaseAggregate
}

func (ar *ITestInterfaceAggregateResult) AggregateResultName() string {
	return "ITestInterface"
}

type ITestInterfaceMutationResult struct {
	ITestInterfaces []*ITestInterface `isResult:"true" json:"iTestInterface,omitempty"`
	NumberAffected  gql.Int           `json:"numUids,omitempty"`
}

type ITestInterfaceDeleteResult struct {
	ITestInterfaceMutationResult
	Message gql.String `json:"msg,omitempty"`
}

func (r ITestInterfaceMutationResult) MutationName() string {
	return "ITestInterface"
}

type ITestInterfaceDeleteArguments struct {
	Filter *ITestInterfaceFilter `json:"filter"`
}

func (t *ITestInterfaceDeleteArguments) IsArg() {}

func (t *ITestInterfaceDeleteArguments) DeleteName() string {
	return "ITestInterfaceFilter"
}

func (t *ITestInterfaceDeleteArguments) DeleteFilter() gql.IFilter {
	return t.Filter
}

type ITestInterfaceFilter struct {
	ID *gql.ID `json:"id,omitempty"`

	Has ITestInterfaceHasChoice `json:"has,omitempty"`
	And []*ITestInterfaceFilter `json:"and,omitempty"`
	Or  []*ITestInterfaceFilter `json:"or,omitempty"`
	Not []*ITestInterfaceFilter `json:"not,omitempty"`
}

func (t *ITestInterfaceFilter) GetName() string {
	return "ITestInterfaceFilter"
}

func (t *ITestInterfaceFilter) IsFilter() {}

func (t *ITestInterfaceFilter) IsArg() {}

type ITestInterfaceOrderChoice string

const ()

type ITestInterfaceOrder struct {
	Asc  ITestInterfaceOrderChoice `json:"asc,omitempty"`
	Desc ITestInterfaceOrderChoice `json:"desc,omitempty"`
	Then *ITestInterfaceOrder      `json:"then,omitempty"`
}

func (t *ITestInterfaceOrder) GetName() string {
	return "ITestInterfaceOrder"
}

func (t *ITestInterfaceOrder) IsOrder() {}

type ITestInterfaceHasChoice string

const (
	HC_ITestInterface_Notes ITestInterfaceHasChoice = "notes"
	HC_ITestInterface_ID    ITestInterfaceHasChoice = "id"
	HC_ITestInterface_Name  ITestInterfaceHasChoice = "name"
)
