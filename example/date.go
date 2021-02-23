package main

import (
	"github.com/machship-oss/go-gql"
)

type Date struct {
	*DateFields

	*DateRelationships
}

func (ar *Date) SingleResultName() string {
	return "Date"
}

type DateMultiResult []*Date

func (ar *DateMultiResult) MultiResultName() string {
	return "Date"
}

type DateFields struct {
	ID             *gql.ID     `json:"id,omitempty"`
	DateCreatedUTC *gql.Time   `json:"dateCreatedUTC,omitempty"`
	Name           *gql.String `json:"name,omitempty"`
}

type DateGetArgs struct {
	ID             *gql.ID     `json:"id,omitempty"`
	DateCreatedUTC *gql.Time   `json:"dateCreatedUTC,omitempty"`
	Name           *gql.String `json:"name,omitempty"`
}

func (t *DateGetArgs) IsArg() {}

func (t *DateGetArgs) IsGet() {}

type DateRelationships struct {
	Year  *Year   `hasInverse:"dates" json:"year,omitempty"`
	Month *Month  `hasInverse:"dates" json:"month,omitempty"`
	Week  *Week   `hasInverse:"dates" json:"week,omitempty"`
	Notes []*Note `hasInverse:"belongsTo" json:"notes,omitempty"`
}

type UpdateDateInput struct {
	Filter *DateFilter `json:"filter"`
	Set    *DatePatch  `json:"set,omitempty"`
	Remove *DatePatch  `json:"remove,omitempty"`
}

func (t *UpdateDateInput) IsArg() {}

func (t *UpdateDateInput) IsUpdate() {}

type DatePatch struct {
	DateFields
	Notes []*gql.BaseRef `json:"notes,omitempty"`
	Year  *gql.BaseRef   `json:"year,omitempty"`
	Month *gql.BaseRef   `json:"month,omitempty"`
	Week  *gql.BaseRef   `json:"week,omitempty"`
}

type AddDateInput struct {
	DateFields
	Notes []*gql.BaseRef `json:"notes,omitempty"`
	Year  *gql.BaseRef   `json:"year,omitempty"`
	Month *gql.BaseRef   `json:"month,omitempty"`
	Week  *gql.BaseRef   `json:"week,omitempty"`
}

type AddDateInputs []*AddDateInput

func (t *AddDateInputs) IsArg() {}

func (t *AddDateInputs) AddName() string {
	return "AddDateInput"
}

type DateAggregateResult struct {
	gql.BaseAggregate
}

func (ar *DateAggregateResult) AggregateResultName() string {
	return "Date"
}

type DateMutationResult struct {
	Dates          []*Date `isResult:"true" json:"date,omitempty"`
	NumberAffected gql.Int `json:"numUids,omitempty"`
}

type DateDeleteResult struct {
	DateMutationResult
	Message gql.String `json:"msg,omitempty"`
}

func (r DateMutationResult) MutationName() string {
	return "Date"
}

type DateDeleteArguments struct {
	Filter *DateFilter `json:"filter"`
}

func (t *DateDeleteArguments) IsArg() {}

func (t *DateDeleteArguments) DeleteName() string {
	return "DateFilter"
}

func (t *DateDeleteArguments) DeleteFilter() gql.IFilter {
	return t.Filter
}

type DateFilter struct {
	ID             *gql.ID   `json:"id,omitempty"`
	DateCreatedUTC *gql.Time `json:"dateCreatedUTC,omitempty"`

	Has DateHasChoice `json:"has,omitempty"`
	And []*DateFilter `json:"and,omitempty"`
	Or  []*DateFilter `json:"or,omitempty"`
	Not []*DateFilter `json:"not,omitempty"`
}

func (t *DateFilter) GetName() string {
	return "DateFilter"
}

func (t *DateFilter) IsFilter() {}

func (t *DateFilter) IsArg() {}

type DateOrderChoice string

const (
	OC_Date_DateCreatedUTC DateOrderChoice = "dateCreatedUTC"
)

type DateOrder struct {
	Asc  DateOrderChoice `json:"asc,omitempty"`
	Desc DateOrderChoice `json:"desc,omitempty"`
	Then *DateOrder      `json:"then,omitempty"`
}

func (t *DateOrder) GetName() string {
	return "DateOrder"
}

func (t *DateOrder) IsOrder() {}

type DateHasChoice string

const (
	HC_Date_Week           DateHasChoice = "week"
	HC_Date_ID             DateHasChoice = "id"
	HC_Date_DateCreatedUTC DateHasChoice = "dateCreatedUTC"
	HC_Date_Name           DateHasChoice = "name"
	HC_Date_Notes          DateHasChoice = "notes"
	HC_Date_Year           DateHasChoice = "year"
	HC_Date_Month          DateHasChoice = "month"
)

func (o *Date) GetWebhookType() string {
	return "com.someone.something.date"
}

func (o *Date) GetWebhookSource() string {
	return "https://github.com/someone/something/v1/type/date"
}

func (o *Date) GetObject() interface{} {
	return o
}
