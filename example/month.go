package main

import (
	"github.com/machship-oss/go-gql"
)

type Month struct {
	*MonthFields

	*MonthRelationships
}

func (ar *Month) SingleResultName() string {
	return "Month"
}

type MonthMultiResult []*Month

func (ar *MonthMultiResult) MultiResultName() string {
	return "Month"
}

type MonthFields struct {
	Name           *gql.String `json:"name,omitempty"`
	Number         *gql.Int64  `json:"number,omitempty"`
	ID             *gql.ID     `json:"id,omitempty"`
	DateCreatedUTC *gql.Time   `json:"dateCreatedUTC,omitempty"`
}

type MonthGetArgs struct {
	Name           *gql.String `json:"name,omitempty"`
	Number         *gql.Int64  `json:"number,omitempty"`
	ID             *gql.ID     `json:"id,omitempty"`
	DateCreatedUTC *gql.Time   `json:"dateCreatedUTC,omitempty"`
}

func (t *MonthGetArgs) IsArg() {}

func (t *MonthGetArgs) IsGet() {}

type MonthRelationships struct {
	Season *Season `hasInverse:"months" json:"season,omitempty"`
	Weeks  []*Week `hasInverse:"month" json:"weeks,omitempty"`
	Dates  []*Date `hasInverse:"month" json:"dates,omitempty"`
	Notes  []*Note `hasInverse:"belongsTo" json:"notes,omitempty"`
}

type UpdateMonthInput struct {
	Filter *MonthFilter `json:"filter"`
	Set    *MonthPatch  `json:"set,omitempty"`
	Remove *MonthPatch  `json:"remove,omitempty"`
}

func (t *UpdateMonthInput) IsArg() {}

func (t *UpdateMonthInput) IsUpdate() {}

type MonthPatch struct {
	MonthFields
	Dates  []*gql.BaseRef `json:"dates,omitempty"`
	Notes  []*gql.BaseRef `json:"notes,omitempty"`
	Season *gql.BaseRef   `json:"season,omitempty"`
	Weeks  []*gql.BaseRef `json:"weeks,omitempty"`
}

type AddMonthInput struct {
	MonthFields
	Season *gql.BaseRef   `json:"season,omitempty"`
	Weeks  []*gql.BaseRef `json:"weeks,omitempty"`
	Dates  []*gql.BaseRef `json:"dates,omitempty"`
	Notes  []*gql.BaseRef `json:"notes,omitempty"`
}

type AddMonthInputs []*AddMonthInput

func (t *AddMonthInputs) IsArg() {}

func (t *AddMonthInputs) AddName() string {
	return "AddMonthInput"
}

type MonthAggregateResult struct {
	gql.BaseAggregate
}

func (ar *MonthAggregateResult) AggregateResultName() string {
	return "Month"
}

type MonthMutationResult struct {
	Months         []*Month `isResult:"true" json:"month,omitempty"`
	NumberAffected gql.Int  `json:"numUids,omitempty"`
}

type MonthDeleteResult struct {
	MonthMutationResult
	Message gql.String `json:"msg,omitempty"`
}

func (r MonthMutationResult) MutationName() string {
	return "Month"
}

type MonthDeleteArguments struct {
	Filter *MonthFilter `json:"filter"`
}

func (t *MonthDeleteArguments) IsArg() {}

func (t *MonthDeleteArguments) DeleteName() string {
	return "MonthFilter"
}

func (t *MonthDeleteArguments) DeleteFilter() gql.IFilter {
	return t.Filter
}

type MonthFilter struct {
	Number         *gql.Int64 `json:"number,omitempty"`
	ID             *gql.ID    `json:"id,omitempty"`
	DateCreatedUTC *gql.Time  `json:"dateCreatedUTC,omitempty"`

	Has MonthHasChoice `json:"has,omitempty"`
	And []*MonthFilter `json:"and,omitempty"`
	Or  []*MonthFilter `json:"or,omitempty"`
	Not []*MonthFilter `json:"not,omitempty"`
}

func (t *MonthFilter) GetName() string {
	return "MonthFilter"
}

func (t *MonthFilter) IsFilter() {}

func (t *MonthFilter) IsArg() {}

type MonthOrderChoice string

const (
	OC_Month_DateCreatedUTC MonthOrderChoice = "dateCreatedUTC"
	OC_Month_Number         MonthOrderChoice = "number"
)

type MonthOrder struct {
	Asc  MonthOrderChoice `json:"asc,omitempty"`
	Desc MonthOrderChoice `json:"desc,omitempty"`
	Then *MonthOrder      `json:"then,omitempty"`
}

func (t *MonthOrder) GetName() string {
	return "MonthOrder"
}

func (t *MonthOrder) IsOrder() {}

type MonthHasChoice string

const (
	HC_Month_Dates          MonthHasChoice = "dates"
	HC_Month_ID             MonthHasChoice = "id"
	HC_Month_DateCreatedUTC MonthHasChoice = "dateCreatedUTC"
	HC_Month_Notes          MonthHasChoice = "notes"
	HC_Month_Name           MonthHasChoice = "name"
	HC_Month_Number         MonthHasChoice = "number"
	HC_Month_Season         MonthHasChoice = "season"
	HC_Month_Weeks          MonthHasChoice = "weeks"
)
