package gql

//Base object:
type testQueryable struct {
	Base
	testQueryableFields
	testQueryableRelationships
}

func (ar *testQueryable) SingleResultName() string {
	return "testQueryable"
}

type testQueryableMultiResult []*testQueryable

func (ar *testQueryableMultiResult) MultiResultName() string {
	return "testQueryable"
}

type testQueryableFields struct {
	TString      *String `json:"tString,omitempty"`
	ShouldIgnore *String `json:"-"`
}

type testQueryableRelationships struct {
	ShouldRecurse *testQueryable `hasInverse:"shouldRecurse" json:"shouldRecurse,omitempty"`
}

//Get:

type testQueryableGet struct {
	ID      *ID     `json:"id,omitempty"`
	TString *String `json:"tString,omitempty"`
}

func (t *testQueryableGet) IsArg() {}
func (t *testQueryableGet) IsGet() {}

//Update:

type updatetestQueryableInput struct {
	Filter *testQueryableFilter `json:"filter"`
	Set    *testQueryablePatch  `json:"set,omitempty"`
	Remove *testQueryablePatch  `json:"remove,omitempty"`
}

func (t *updatetestQueryableInput) IsArg()    {}
func (t *updatetestQueryableInput) IsUpdate() {}

type testQueryablePatch struct {
	BasePatch
	testQueryableFields
	ShouldRecurse []*BaseRef `hasInverse:"shouldRecurse" json:"shouldRecurse,omitempty"`
}

//Add:

type AddtestQueryableInput struct {
	AddBaseInput
	testQueryableFields
	SomeRef []*BaseRef `hasInverse:"someRef" json:"someRef,omitempty"`
}

type AddtestQueryableInputs []*AddtestQueryableInput

func (t *AddtestQueryableInputs) IsArg() {}
func (t *AddtestQueryableInputs) AddName() string {
	return "AddtestQueryableInput"
}

//Aggregate:

type testQueryableAggregateResult struct {
	BaseAggregate
}

func (ar *testQueryableAggregateResult) AggregateResultName() string {
	return "testQueryable"
}

//Mutation Result:

type testQueryableMutationResult struct {
	TestQueryables []*testQueryable `isResult:"true" json:"testQueryable,omitempty"`
	NumberAffected Int              `json:"numUids,omitempty"`
}

type testQueryableDeleteResult struct {
	testQueryableMutationResult
	Message String `json:"msg,omitempty"`
}

func (r testQueryableMutationResult) MutationName() string {
	return "testQueryable"
}

//Filters:

type testQueryableDeleteArguments struct {
	Filter *testQueryableFilter `json:"filter"`
}

func (t *testQueryableDeleteArguments) IsArg() {}

func (t *testQueryableDeleteArguments) DeleteName() string {
	return "testQueryableFilter"
}

func (t *testQueryableDeleteArguments) DeleteFilter() IFilter {
	return t.Filter
}

//Filter:
type testQueryableFilter struct {
	BaseFilter
	TString      *StringFilter `json:"tString,omitempty"`
	ShouldIgnore *StringFilter `json:"-"`

	Has testQueryableHasChoice `json:"has,omitempty"`
	And []testQueryableFilter  `json:"and,omitempty"`
	Or  []testQueryableFilter  `json:"or,omitempty"`
	Not []testQueryableFilter  `json:"not,omitempty"`
}

func (t *testQueryableFilter) IsArg() {}
func (t *testQueryableFilter) GetName() string {
	return "testQueryableFilter"
}

type testQueryableHasChoice string

const (
	HC_testQueryable_ID             testQueryableHasChoice = "id"
	HC_testQueryable_DateCreatedUTC testQueryableHasChoice = "dateCreatedUTC"
	HC_testQueryable_Name           testQueryableHasChoice = "name"
)

func (t *testQueryableFilter) IsFilter() {}

//Order:

type testQueryableOrderChoice BaseOrderChoice

const (
	OC_SomeFakeField testQueryableOrderChoice = "someFakeField"
)

type testQueryableOrder struct {
	Asc  testQueryableOrderChoice `json:"asc,omitempty"`
	Desc testQueryableOrderChoice `json:"desc,omitempty"`
	Then *testQueryableOrder      `json:"then,omitempty"`
}

func (t *testQueryableOrder) GetName() string {
	return "testQueryableOrder"
}

func (t *testQueryableOrder) IsOrder() {}

type shouldForceJSONError struct {
	FailMe chan int `json:"bad,omitempty"`
}

func (t *shouldForceJSONError) GetName() string {
	return "shouldFail"
}

func (t *shouldForceJSONError) IsFilter() {}
func (t *shouldForceJSONError) IsArg()    {}
