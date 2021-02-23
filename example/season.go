package main

import (
	"github.com/machship-oss/go-gql"
)

type Season struct {
	*SeasonFields

	*SeasonRelationships
}

func (ar *Season) SingleResultName() string {
	return "Season"
}

type SeasonMultiResult []*Season

func (ar *SeasonMultiResult) MultiResultName() string {
	return "Season"
}

type SeasonFields struct {
	ID             *gql.ID     `json:"id,omitempty"`
	DateCreatedUTC *gql.Time   `json:"dateCreatedUTC,omitempty"`
	Name           *gql.String `json:"name,omitempty"`
}

type SeasonGetArgs struct {
	ID             *gql.ID     `json:"id,omitempty"`
	DateCreatedUTC *gql.Time   `json:"dateCreatedUTC,omitempty"`
	Name           *gql.String `json:"name,omitempty"`
}

func (t *SeasonGetArgs) IsArg() {}

func (t *SeasonGetArgs) IsGet() {}

type SeasonRelationships struct {
	Notes  []*Note  `hasInverse:"belongsTo" json:"notes,omitempty"`
	Months []*Month `hasInverse:"season" json:"months,omitempty"`
}

type UpdateSeasonInput struct {
	Filter *SeasonFilter `json:"filter"`
	Set    *SeasonPatch  `json:"set,omitempty"`
	Remove *SeasonPatch  `json:"remove,omitempty"`
}

func (t *UpdateSeasonInput) IsArg() {}

func (t *UpdateSeasonInput) IsUpdate() {}

type SeasonPatch struct {
	SeasonFields
	Notes  []*gql.BaseRef `json:"notes,omitempty"`
	Months []*gql.BaseRef `json:"months,omitempty"`
}

type AddSeasonInput struct {
	SeasonFields
	Notes  []*gql.BaseRef `json:"notes,omitempty"`
	Months []*gql.BaseRef `json:"months,omitempty"`
}

type AddSeasonInputs []*AddSeasonInput

func (t *AddSeasonInputs) IsArg() {}

func (t *AddSeasonInputs) AddName() string {
	return "AddSeasonInput"
}

type SeasonAggregateResult struct {
	gql.BaseAggregate
}

func (ar *SeasonAggregateResult) AggregateResultName() string {
	return "Season"
}

type SeasonMutationResult struct {
	Seasons        []*Season `isResult:"true" json:"season,omitempty"`
	NumberAffected gql.Int   `json:"numUids,omitempty"`
}

type SeasonDeleteResult struct {
	SeasonMutationResult
	Message gql.String `json:"msg,omitempty"`
}

func (r SeasonMutationResult) MutationName() string {
	return "Season"
}

type SeasonDeleteArguments struct {
	Filter *SeasonFilter `json:"filter"`
}

func (t *SeasonDeleteArguments) IsArg() {}

func (t *SeasonDeleteArguments) DeleteName() string {
	return "SeasonFilter"
}

func (t *SeasonDeleteArguments) DeleteFilter() gql.IFilter {
	return t.Filter
}

type SeasonFilter struct {
	ID             *gql.ID   `json:"id,omitempty"`
	DateCreatedUTC *gql.Time `json:"dateCreatedUTC,omitempty"`

	Has SeasonHasChoice `json:"has,omitempty"`
	And []*SeasonFilter `json:"and,omitempty"`
	Or  []*SeasonFilter `json:"or,omitempty"`
	Not []*SeasonFilter `json:"not,omitempty"`
}

func (t *SeasonFilter) GetName() string {
	return "SeasonFilter"
}

func (t *SeasonFilter) IsFilter() {}

func (t *SeasonFilter) IsArg() {}

type SeasonOrderChoice string

const (
	OC_Season_DateCreatedUTC SeasonOrderChoice = "dateCreatedUTC"
)

type SeasonOrder struct {
	Asc  SeasonOrderChoice `json:"asc,omitempty"`
	Desc SeasonOrderChoice `json:"desc,omitempty"`
	Then *SeasonOrder      `json:"then,omitempty"`
}

func (t *SeasonOrder) GetName() string {
	return "SeasonOrder"
}

func (t *SeasonOrder) IsOrder() {}

type SeasonHasChoice string

const (
	HC_Season_Months         SeasonHasChoice = "months"
	HC_Season_ID             SeasonHasChoice = "id"
	HC_Season_DateCreatedUTC SeasonHasChoice = "dateCreatedUTC"
	HC_Season_Name           SeasonHasChoice = "name"
	HC_Season_Notes          SeasonHasChoice = "notes"
)
