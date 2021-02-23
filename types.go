package gql

import (
	"encoding/json"
	"fmt"
	"time"
)

type Bool struct {
	Bool  bool
	Valid bool // Valid is true if Bool is not NULL
}

func NewBoolStruct(x bool) Bool {
	nw := NewBool(x)
	return *nw
}

func NewBool(x bool) *Bool { return &Bool{Bool: x, Valid: true} }

func (nb Bool) GetName() string {
	return "Boolean"
}

func (nb Bool) MarshalJSON() ([]byte, error) {
	if nb.Valid {
		return json.Marshal(nb.Bool)
	}
	return json.Marshal(nil)
}

func (nb *Bool) UnmarshalJSON(data []byte) error {
	var b *bool
	if err := json.Unmarshal(data, &b); err != nil {
		return err
	}
	if b != nil {
		nb.Valid = true
		nb.Bool = *b
	} else {
		nb.Valid = false
	}
	return nil
}

type Float64 struct {
	Float64 float64
	Valid   bool // Valid is true if Float64 is not NULL
}

func NewFloat64Struct(x float64) Float64 {
	nw := NewFloat64(x)
	return *nw
}

func NewFloat64(x float64) *Float64 { return &Float64{Float64: x, Valid: true} }

func (nb Float64) GetName() string {
	return "Float"
}

func (nf Float64) MarshalJSON() ([]byte, error) {
	if nf.Valid {
		return json.Marshal(nf.Float64)
	}
	return json.Marshal(nil)
}

func (nf *Float64) UnmarshalJSON(data []byte) error {
	var f *float64
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}
	if f != nil {
		nf.Valid = true
		nf.Float64 = *f
	} else {
		nf.Valid = false
	}
	return nil
}

type Int64 struct {
	Int64 int64
	Valid bool // Valid is true if Int64 is not NULL
}

func NewInt64Struct(x int64) Int64 {
	nw := NewInt64(x)
	return *nw
}

func NewInt64(x int64) *Int64 { return &Int64{Int64: x, Valid: true} }

func (nb Int64) GetName() string {
	return "Int64"
}

func (ni Int64) MarshalJSON() ([]byte, error) {
	if ni.Valid {
		return json.Marshal(ni.Int64)
	}
	return json.Marshal(nil)
}

func (ni *Int64) UnmarshalJSON(data []byte) error {
	var i *int64
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	if i != nil {
		ni.Valid = true
		ni.Int64 = *i
	} else {
		ni.Valid = false
	}
	return nil
}

type Int struct {
	Int   int
	Valid bool // Valid is true if Int is not NULL
}

func NewIntStruct(x int) Int {
	nw := NewInt(x)
	return *nw
}

func NewInt(x int) *Int { return &Int{Int: x, Valid: true} }

func (nb Int) GetName() string {
	return "Int"
}

func (ni Int) MarshalJSON() ([]byte, error) {
	if ni.Valid {
		return json.Marshal(ni.Int)
	}
	return json.Marshal(nil)
}

func (ni *Int) UnmarshalJSON(data []byte) error {
	var i *int
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	if i != nil {
		ni.Valid = true
		ni.Int = *i
	} else {
		ni.Valid = false
	}
	return nil
}

type String struct {
	String string
	Valid  bool // Valid is true if String is not NULL
}

func NewStringStruct(x string) String {
	nw := NewString(x)
	return *nw
}

func NewString(x string) *String { return &String{String: x, Valid: true} }

func (nb String) GetName() string {
	return "String"
}

func (ns String) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return json.Marshal(nil)
}

func (ns *String) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		ns.Valid = true
		ns.String = *s
	} else {
		ns.Valid = false
	}
	return nil
}

type Time struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

func NewTimeStruct(x time.Time) Time {
	nw := NewTime(x)
	return *nw
}

func NewTime(x time.Time) *Time { return &Time{Time: x, Valid: true} }

func (nb Time) GetName() string {
	return "DateTime"
}

func (nt Time) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		return json.Marshal(nt.Time)
	}
	return json.Marshal(nil)
}

func (nt *Time) UnmarshalJSON(data []byte) error {
	var t *time.Time
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	if t != nil {
		nt.Valid = true
		nt.Time = *t
	} else {
		nt.Valid = false
	}
	return nil
}

type Point struct {
	//this does not need an ID as it is a scalar type in GraphQL

	Latitude  *Float64 `json:"latitude,omitempty"`
	Longitude *Float64 `json:"longitude,omitempty"`
}

func NewPointStruct(lat, lng float64) Point {
	nw := NewPoint(lat, lng)
	return *nw
}

func NewPoint(lat, lng float64) *Point {
	return &Point{Latitude: NewFloat64(lat), Longitude: NewFloat64(lng)}
}

func (nb Point) GetName() string {
	return "Point"
}

type ID struct {
	ID    string
	Valid bool // Valid is true if ID is not NULL
}

func NewIDStruct(x string) ID {
	nw := NewID(x)
	return *nw
}

func NewID(x string) *ID { return &ID{ID: x, Valid: true} }

func (nb ID) GetName() string {
	return "ID"
}

func (ns ID) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.ID)
	}
	return json.Marshal(nil)
}

func (ns *ID) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		ns.Valid = true
		ns.ID = *s
	} else {
		ns.Valid = false
	}
	return nil
}

func newErr(outer, inner error) GqlError {
	return GqlError{
		error:      outer,
		InnerError: inner,
	}
}

func (e GqlError) String() string {
	return fmt.Sprintf("%v: %v", e.error.Error(), e.InnerError.Error())
}

type GqlError struct {
	error
	InnerError error
}

type IOrder interface {
	IsOrder()
	GetName() string
}

type IFilter interface {
	IsFilter()
	IsArg()
	GetName() string
}

type IIsArg interface {
	IsArg()
}

type IIsGet interface {
	IsGet()
	IsArg()
}

type IIsAdd interface {
	AddName() string
	IsArg()
}

type IIsDelete interface {
	DeleteName() string
	DeleteFilter() IFilter
	IsArg()
}

type IIsUpdate interface {
	IsUpdate()
	IsArg()
}

type IMutationResult interface {
	MutationName() string
}

type ISingleResult interface {
	SingleResultName() string
}

type IMultiResult interface {
	MultiResultName() string
}

type IAggregateResult interface {
	AggregateResultName() string
}

type argType uint8

const (
	at_Get argType = iota
	at_Add
	at_Update
	at_Delete
)

// Base types:

type Base struct {
	ID             *ID   `json:"id,omitempty"`
	DateCreatedUTC *Time `json:"dateCreatedUTC,omitempty"`

	//TypeName is readonly and can only be queried from the server
	TypeName *String `json:"__typename,omitempty"`
}

type BaseGet struct {
	ID             *ID   `json:"id,omitempty"`
	DateCreatedUTC *Time `json:"dateCreatedUTC,omitempty"`
}

type BaseFields struct {
}

type BaseFilter struct {
	ID          *ID             `json:"id,omitempty"`
	DateCreated *DateTimeFilter `json:"dateCreated,omitempty"`
}

type DateTimeRange struct {
	Min *Time `json:"min"`
	Max *Time `json:"max"`
}

type BaseRef struct {
	ID *ID `json:"id,omitempty"`
}

type StringFilter struct {
	Equal *String `json:"eq,omitempty"`
}

type DateTimeFilter struct {
	Equal              *Time          `json:"eq,omitempty"`
	LessThanOrEqual    *Time          `json:"le"`
	LessThan           *Time          `json:"lt"`
	GreaterThanOrEqual *Time          `json:"ge"`
	GreaterThan        *Time          `json:"gt"`
	Between            *DateTimeRange `json:"between"`
}

type BaseOrderChoice string

const (
	OC_BaseDateCreated BaseOrderChoice = "dateCreated"
)

func (c *BaseOrderChoice) GetName() string {
	return string(*c)
}

type BaseHasChoice string

const (
	HC_BaseDateCreated BaseHasChoice = "dateCreated"
)

func (c *BaseHasChoice) GetName() string {
	return string(*c)
}

type AddBaseInput struct {
	DateCreated *Time `json:"dateCreated,omitempty"`
}

type BasePatch struct {
	DateCreated *Time `json:"dateCreated,omitempty"`
}

type BaseAggregate struct {
	Count int `json:"count,omitempty"`
}
