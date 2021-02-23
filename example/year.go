package main

import (
	"github.com/machship-oss/go-gql"
)

type Year struct {
	*YearFields

	*YearRelationships
}

func (ar *Year) SingleResultName() string {
	return "Year"
}

type YearMultiResult []*Year

func (ar *YearMultiResult) MultiResultName() string {
	return "Year"
}

type YearFields struct {
	ID             *gql.ID    `json:"id,omitempty"`
	DateCreatedUTC *gql.Time  `json:"dateCreatedUTC,omitempty"`
	Year           *gql.Int64 `json:"year,omitempty"`
}

type YearGetArgs struct {
	Year           *gql.Int64 `json:"year,omitempty"`
	ID             *gql.ID    `json:"id,omitempty"`
	DateCreatedUTC *gql.Time  `json:"dateCreatedUTC,omitempty"`
}

func (t *YearGetArgs) IsArg() {}

func (t *YearGetArgs) IsGet() {}

type YearRelationships struct {
	Dates []*Date `hasInverse:"year" json:"dates,omitempty"`
	Weeks []*Week `hasInverse:"year" json:"weeks,omitempty"`
}

type UpdateYearInput struct {
	Filter *YearFilter `json:"filter"`
	Set    *YearPatch  `json:"set,omitempty"`
	Remove *YearPatch  `json:"remove,omitempty"`
}

func (t *UpdateYearInput) IsArg() {}

func (t *UpdateYearInput) IsUpdate() {}

type YearPatch struct {
	YearFields
	Dates        []*gql.BaseRef `json:"dates,omitempty"`
	Weeks        []*gql.BaseRef `json:"weeks,omitempty"`
	FakePassword *gql.String    `json:"fakePassword,omitempty"`
}

type AddYearInput struct {
	YearFields
	FakePassword *gql.String    `json:"fakePassword,omitempty"`
	Dates        []*gql.BaseRef `json:"dates,omitempty"`
	Weeks        []*gql.BaseRef `json:"weeks,omitempty"`
}

type AddYearInputs []*AddYearInput

func (t *AddYearInputs) IsArg() {}

func (t *AddYearInputs) AddName() string {
	return "AddYearInput"
}

type YearAggregateResult struct {
	gql.BaseAggregate
}

func (ar *YearAggregateResult) AggregateResultName() string {
	return "Year"
}

type YearMutationResult struct {
	Years          []*Year `isResult:"true" json:"year,omitempty"`
	NumberAffected gql.Int `json:"numUids,omitempty"`
}

type YearDeleteResult struct {
	YearMutationResult
	Message gql.String `json:"msg,omitempty"`
}

func (r YearMutationResult) MutationName() string {
	return "Year"
}

type YearDeleteArguments struct {
	Filter *YearFilter `json:"filter"`
}

func (t *YearDeleteArguments) IsArg() {}

func (t *YearDeleteArguments) DeleteName() string {
	return "YearFilter"
}

func (t *YearDeleteArguments) DeleteFilter() gql.IFilter {
	return t.Filter
}

type YearFilter struct {
	ID             *gql.ID   `json:"id,omitempty"`
	DateCreatedUTC *gql.Time `json:"dateCreatedUTC,omitempty"`

	Has YearHasChoice `json:"has,omitempty"`
	And []*YearFilter `json:"and,omitempty"`
	Or  []*YearFilter `json:"or,omitempty"`
	Not []*YearFilter `json:"not,omitempty"`
}

func (t *YearFilter) GetName() string {
	return "YearFilter"
}

func (t *YearFilter) IsFilter() {}

func (t *YearFilter) IsArg() {}

type YearOrderChoice string

const (
	OC_Year_DateCreatedUTC YearOrderChoice = "dateCreatedUTC"
)

type YearOrder struct {
	Asc  YearOrderChoice `json:"asc,omitempty"`
	Desc YearOrderChoice `json:"desc,omitempty"`
	Then *YearOrder      `json:"then,omitempty"`
}

func (t *YearOrder) GetName() string {
	return "YearOrder"
}

func (t *YearOrder) IsOrder() {}

type YearHasChoice string

const (
	HC_Year_Year           YearHasChoice = "year"
	HC_Year_Dates          YearHasChoice = "dates"
	HC_Year_Weeks          YearHasChoice = "weeks"
	HC_Year_ID             YearHasChoice = "id"
	HC_Year_DateCreatedUTC YearHasChoice = "dateCreatedUTC"
)

func (o *Year) CheckPassword(Year *gql.String) (isCorrect bool) {
	//do some password stuff
	return false
}
