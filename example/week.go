package main

import (
	"github.com/machship-oss/go-gql"
)

type Week struct {
	*WeekFields

	*WeekRelationships
}

func (ar *Week) SingleResultName() string {
	return "Week"
}

type WeekMultiResult []*Week

func (ar *WeekMultiResult) MultiResultName() string {
	return "Week"
}

type WeekFields struct {
	NumberInYear   *gql.Int64 `json:"numberInYear,omitempty"`
	ID             *gql.ID    `json:"id,omitempty"`
	DateCreatedUTC *gql.Time  `json:"dateCreatedUTC,omitempty"`
}

type WeekGetArgs struct {
	NumberInYear   *gql.Int64 `json:"numberInYear,omitempty"`
	ID             *gql.ID    `json:"id,omitempty"`
	DateCreatedUTC *gql.Time  `json:"dateCreatedUTC,omitempty"`
}

func (t *WeekGetArgs) IsArg() {}

func (t *WeekGetArgs) IsGet() {}

type WeekRelationships struct {
	Year  *Year   `hasInverse:"weeks" json:"year,omitempty"`
	Month *Month  `hasInverse:"weeks" json:"month,omitempty"`
	Dates []*Date `hasInverse:"week" json:"dates,omitempty"`
}

type UpdateWeekInput struct {
	Filter *WeekFilter `json:"filter"`
	Set    *WeekPatch  `json:"set,omitempty"`
	Remove *WeekPatch  `json:"remove,omitempty"`
}

func (t *UpdateWeekInput) IsArg() {}

func (t *UpdateWeekInput) IsUpdate() {}

type WeekPatch struct {
	WeekFields
	Year  *gql.BaseRef   `json:"year,omitempty"`
	Month *gql.BaseRef   `json:"month,omitempty"`
	Dates []*gql.BaseRef `json:"dates,omitempty"`
}

type AddWeekInput struct {
	WeekFields
	Month *gql.BaseRef   `json:"month,omitempty"`
	Dates []*gql.BaseRef `json:"dates,omitempty"`
	Year  *gql.BaseRef   `json:"year,omitempty"`
}

type AddWeekInputs []*AddWeekInput

func (t *AddWeekInputs) IsArg() {}

func (t *AddWeekInputs) AddName() string {
	return "AddWeekInput"
}

type WeekAggregateResult struct {
	gql.BaseAggregate
}

func (ar *WeekAggregateResult) AggregateResultName() string {
	return "Week"
}

type WeekMutationResult struct {
	Weeks          []*Week `isResult:"true" json:"week,omitempty"`
	NumberAffected gql.Int `json:"numUids,omitempty"`
}

type WeekDeleteResult struct {
	WeekMutationResult
	Message gql.String `json:"msg,omitempty"`
}

func (r WeekMutationResult) MutationName() string {
	return "Week"
}

type WeekDeleteArguments struct {
	Filter *WeekFilter `json:"filter"`
}

func (t *WeekDeleteArguments) IsArg() {}

func (t *WeekDeleteArguments) DeleteName() string {
	return "WeekFilter"
}

func (t *WeekDeleteArguments) DeleteFilter() gql.IFilter {
	return t.Filter
}

type WeekFilter struct {
	DateCreatedUTC *gql.Time `json:"dateCreatedUTC,omitempty"`
	ID             *gql.ID   `json:"id,omitempty"`

	Has WeekHasChoice `json:"has,omitempty"`
	And []*WeekFilter `json:"and,omitempty"`
	Or  []*WeekFilter `json:"or,omitempty"`
	Not []*WeekFilter `json:"not,omitempty"`
}

func (t *WeekFilter) GetName() string {
	return "WeekFilter"
}

func (t *WeekFilter) IsFilter() {}

func (t *WeekFilter) IsArg() {}

type WeekOrderChoice string

const (
	OC_Week_DateCreatedUTC WeekOrderChoice = "dateCreatedUTC"
)

type WeekOrder struct {
	Asc  WeekOrderChoice `json:"asc,omitempty"`
	Desc WeekOrderChoice `json:"desc,omitempty"`
	Then *WeekOrder      `json:"then,omitempty"`
}

func (t *WeekOrder) GetName() string {
	return "WeekOrder"
}

func (t *WeekOrder) IsOrder() {}

type WeekHasChoice string

const (
	HC_Week_NumberInYear   WeekHasChoice = "numberInYear"
	HC_Week_Year           WeekHasChoice = "year"
	HC_Week_Month          WeekHasChoice = "month"
	HC_Week_Dates          WeekHasChoice = "dates"
	HC_Week_ID             WeekHasChoice = "id"
	HC_Week_DateCreatedUTC WeekHasChoice = "dateCreatedUTC"
)
