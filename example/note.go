package main

import (
	"github.com/machship-oss/go-gql"
)

type Note struct {
	*NoteFields

	*NoteRelationships
}

func (ar *Note) SingleResultName() string {
	return "Note"
}

type NoteMultiResult []*Note

func (ar *NoteMultiResult) MultiResultName() string {
	return "Note"
}

type NoteFields struct {
	ID             *gql.ID     `json:"id,omitempty"`
	DateCreatedUTC *gql.Time   `json:"dateCreatedUTC,omitempty"`
	Note           *gql.String `json:"note,omitempty"`
}

type NoteGetArgs struct {
	Note           *gql.String `json:"note,omitempty"`
	ID             *gql.ID     `json:"id,omitempty"`
	DateCreatedUTC *gql.Time   `json:"dateCreatedUTC,omitempty"`
}

func (t *NoteGetArgs) IsArg() {}

func (t *NoteGetArgs) IsGet() {}

type NoteRelationships struct {
	BelongsTo *ITestInterface `hasInverse:"notes" json:"belongsTo,omitempty"`
}

type UpdateNoteInput struct {
	Filter *NoteFilter `json:"filter"`
	Set    *NotePatch  `json:"set,omitempty"`
	Remove *NotePatch  `json:"remove,omitempty"`
}

func (t *UpdateNoteInput) IsArg() {}

func (t *UpdateNoteInput) IsUpdate() {}

type NotePatch struct {
	NoteFields
	BelongsTo *gql.BaseRef `json:"belongsTo,omitempty"`
}

type AddNoteInput struct {
	NoteFields
	BelongsTo *gql.BaseRef `json:"belongsTo,omitempty"`
}

type AddNoteInputs []*AddNoteInput

func (t *AddNoteInputs) IsArg() {}

func (t *AddNoteInputs) AddName() string {
	return "AddNoteInput"
}

type NoteAggregateResult struct {
	gql.BaseAggregate
}

func (ar *NoteAggregateResult) AggregateResultName() string {
	return "Note"
}

type NoteMutationResult struct {
	Notes          []*Note `isResult:"true" json:"note,omitempty"`
	NumberAffected gql.Int `json:"numUids,omitempty"`
}

type NoteDeleteResult struct {
	NoteMutationResult
	Message gql.String `json:"msg,omitempty"`
}

func (r NoteMutationResult) MutationName() string {
	return "Note"
}

type NoteDeleteArguments struct {
	Filter *NoteFilter `json:"filter"`
}

func (t *NoteDeleteArguments) IsArg() {}

func (t *NoteDeleteArguments) DeleteName() string {
	return "NoteFilter"
}

func (t *NoteDeleteArguments) DeleteFilter() gql.IFilter {
	return t.Filter
}

type NoteFilter struct {
	ID             *gql.ID   `json:"id,omitempty"`
	DateCreatedUTC *gql.Time `json:"dateCreatedUTC,omitempty"`

	Has NoteHasChoice `json:"has,omitempty"`
	And []*NoteFilter `json:"and,omitempty"`
	Or  []*NoteFilter `json:"or,omitempty"`
	Not []*NoteFilter `json:"not,omitempty"`
}

func (t *NoteFilter) GetName() string {
	return "NoteFilter"
}

func (t *NoteFilter) IsFilter() {}

func (t *NoteFilter) IsArg() {}

type NoteOrderChoice string

const (
	OC_Note_DateCreatedUTC NoteOrderChoice = "dateCreatedUTC"
)

type NoteOrder struct {
	Asc  NoteOrderChoice `json:"asc,omitempty"`
	Desc NoteOrderChoice `json:"desc,omitempty"`
	Then *NoteOrder      `json:"then,omitempty"`
}

func (t *NoteOrder) GetName() string {
	return "NoteOrder"
}

func (t *NoteOrder) IsOrder() {}

type NoteHasChoice string

const (
	HC_Note_ID             NoteHasChoice = "id"
	HC_Note_DateCreatedUTC NoteHasChoice = "dateCreatedUTC"
	HC_Note_Note           NoteHasChoice = "note"
	HC_Note_BelongsTo      NoteHasChoice = "belongsTo"
)
