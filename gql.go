package gql

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync/atomic"
	"unicode"
)

type fieldMap struct {
	isTop     bool
	typeName  string
	fields    []string
	subFields []*fieldMap
}

func (fm *fieldMap) writeString(depth int) string {
	if len(fm.fields) == 0 && len(fm.subFields) == 0 {
		return ""
	}
	sb := &strings.Builder{}

	if !fm.isTop {
		sb.WriteString(fm.typeName)
	}

	sb.WriteString("{")
	sb.WriteString(strings.Join(fm.fields, ","))
	if len(fm.subFields) > 0 {
		sfStrings := []string{}
		for _, sf := range fm.subFields {
			sfs := sf.writeString(depth + 1)
			if sfs != "" {
				sfStrings = append(sfStrings, sfs)
			}
		}
		if len(sfStrings) > 0 {
			if len(fm.fields) > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(strings.Join(sfStrings, ","))
		}
	}
	sb.WriteString("}")

	return sb.String()
}

func writeQuery2(t reflect.Type, fm *fieldMap, seen map[string]struct{}, wasAnonymous bool, nextName string) (wasStruct bool) {
	if seen == nil {
		seen = make(map[string]struct{})
	}
	var thisFM *fieldMap

	subFields := []*subField{}

	switch k := t.Kind(); k {
	case reflect.Ptr, reflect.Slice, reflect.Array:
		// For Pointers, we need to get the concrete type, so we go to its element
		// Slices and Arrays are treated the same as a single sub object by GraphQL, so just go to its element
		te := t.Elem()
		return writeQuery2(te, fm, seen, wasAnonymous, nextName)
	case reflect.Struct:
		wasStruct = true
		// If the type implements json.Unmarshaler, we treat is as a Scalar value and just take it's JSON tag so
		// we don't expand it's fields (e.g. gql.String); when deserialising, json will treat it as a flat field.
		if reflect.PtrTo(t).Implements(jsonUnmarshaler) {
			wasStruct = false
			return
		}

		if fm.typeName == "" {
			fm.typeName = nextName
			thisFM = fm
		} else if wasAnonymous {
			thisFM = fm
		} else {
			thisFM = &fieldMap{
				typeName: nextName,
			}
			fm.subFields = append(fm.subFields, thisFM)
		}

		// Deal with the fields
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)

			// Anonymous means that it is an embedded field, so we don't write it's name
			// if f.Anonymous {
			// 	writeQuery2(f.Type, thisFM, seen, true, nextName)
			// 	continue
			// }

			fieldName, _ := f.Tag.Lookup("json")
			fieldName = strings.Replace(fieldName, ",omitempty", "", 1)
			switch fieldName {
			case "-":
				// Skip this field as it is meant to be skipped by JSON
				continue
			case "":
				fieldName = lowerFirst(f.Name)
			}

			var thisInverseName string
			if !f.Anonymous {
				thisInverseName, _ = f.Tag.Lookup("hasInverse")
			}

			myName := fmt.Sprintf("%s_%s_%s", nextName, f.Name, thisInverseName)
			if _, ok := seen[myName]; ok {
				//Don't allow recursion any further
				continue
			}
			seen[myName] = struct{}{}

			if fieldName == "shouldRecurse" {
				fmt.Print()
			}

			subFields = append(subFields, &subField{fieldName, f, f.Anonymous})
		}
	}

	for _, sf := range subFields {
		var nm string
		if !sf.WasAnonymous {
			nm = sf.Name
		} else {
			nm = thisFM.typeName
		}
		subWasStruct := writeQuery2(sf.Field.Type, thisFM, seen, sf.WasAnonymous, nm)
		if !subWasStruct && !sf.WasAnonymous {
			thisFM.fields = append(thisFM.fields, sf.Name)
		}
	}

	return wasStruct
}

type subField struct {
	Name         string
	Field        reflect.StructField
	WasAnonymous bool
}

var jsonUnmarshaler = reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()

type queryPart struct {
	identifier       string
	query            string
	innerQueryParams []string
	outerQueryParams []string
	argumentsJSON    []string
	argumentJSON     string

	queries, queryParams, queryVariablesJSON []string
	//todo: clean up these names so that they have more meaning
}

var counter = int64(-1)
var testMode bool

func (op *queryPart) prepareQuery(opName, resultName string, isTestMode ...bool) (qsb *strings.Builder) {
	thisTestMode := testMode
	if len(isTestMode) == 1 {
		thisTestMode = isTestMode[0]
	}
	if op.identifier == "" {
		if thisTestMode {
			op.identifier = fmt.Sprint(0)
		} else {
			cntrVal := atomic.AddInt64(&counter, 1)
			op.identifier = fmt.Sprint(cntrVal)
		}
	}

	op.identifier = fmt.Sprintf("%s_%s", opName, op.identifier)

	qsb = &strings.Builder{}

	// write the query identifier
	qsb.WriteString(fmt.Sprintf("%s:", op.identifier))
	qsb.WriteString(opName + resultName)
	return qsb
}

func (op *queryPart) finaliseQuery(qsb *strings.Builder, resultObject interface{}) {
	if len(op.innerQueryParams) > 0 {
		qsb.WriteString(fmt.Sprintf("(%s)", strings.Join(op.innerQueryParams, ",")))
	}

	// var fsb strings.Builder
	// writeQuery(&fsb, reflect.TypeOf(resultObject), false, nil)

	fm := fieldMap{isTop: true}
	writeQuery2(reflect.TypeOf(resultObject), &fm, nil, false, "top")
	qry := fm.writeString(0)

	// qsb.WriteString(fsb.String())
	qsb.WriteString(qry)
	qrtString := qsb.String()
	op.query = qrtString
}

func (op *queryPart) addParamFieldsAndJson(thing IIsArg, argtyp argType, argName ...string) (err error) {
	// input could be a pointer
	rStructValue := reflect.ValueOf(thing)
	if rStructValue.Kind() == reflect.Ptr {
		rStructValue = rStructValue.Elem()
	}

	// we make sure it is a struct:
	if !(rStructValue.Kind() == reflect.Struct || rStructValue.Kind() == reflect.Array || rStructValue.Kind() == reflect.Slice) {
		return ErrArgumentWasNotAStruct
	}

	switch argtyp {
	case at_Add:
		if len(argName) == 0 {
			return newErr(ErrArgumentError, err)
		}

		nm := fmt.Sprintf("[%s!]!", argName[0])
		err = op.addParamsAndJson(thing, nm, "input")
		if err != nil {
			return newErr(ErrArgumentError, err)
		}
	case at_Update:
		nm := fmt.Sprintf("%s!", rStructValue.Type().Name())
		err = op.addParamsAndJson(thing, nm, "input")
		if err != nil {
			return newErr(ErrArgumentError, err)
		}
	case at_Delete:
		if len(argName) == 0 {
			return newErr(ErrArgumentError, err)
		}

		nm := fmt.Sprintf("%s!", argName[0])
		err = op.addParamsAndJson(thing, nm, "filter")
		if err != nil {
			return newErr(ErrArgumentError, err)
		}
	case at_Get:
		rStructType := rStructValue.Type()
		// we now traverse the fields
		for i := 0; i < rStructType.NumField(); i++ {
			// var sb strings.Builder
			f := rStructType.Field(i)
			var fldName string
			value, _ := f.Tag.Lookup("json")
			switch value {
			case "-":
				// Skip this field as it is meant to be skipped by JSON
				continue
			case "":
				value = lowerFirst(f.Name)
			}

			fldVal := rStructValue.FieldByIndex(f.Index).Interface()

			var typName string
			fldName = strings.Replace(value, ",omitempty", "", 1)
			if f.Type.Kind() == reflect.Struct {
				typName = f.Type.Name()
			} else {
				typName = f.Type.Elem().Name()
				//check if field value is nil. If so, omit it
				if fld := rStructValue.Field(i); fld.IsNil() || (fld.Kind() == reflect.Struct && fld.IsZero()) {
					//ignore this field
					continue
				}
			}

			// writeArgumentType(&sb, f.Type, true, false)
			err = op.addParamsAndJson(fldVal, typName, fldName)
			if err != nil {
				return newErr(ErrArgumentError, err)
			}
		}
	}

	return nil
}

func (op *queryPart) addParamsAndJson(thing interface{}, argName, jsonName string) (err error) {
	op.outerQueryParams = append(op.outerQueryParams, fmt.Sprintf("$%s_%s:%s", jsonName, op.identifier, argName))
	op.innerQueryParams = append(op.innerQueryParams, fmt.Sprintf("%s:$%s_%s", jsonName, jsonName, op.identifier))
	argJSON, err := json.Marshal(thing)
	if err != nil {
		return newErr(ErrArgumentCouldNotBeEncodedAsJSON, err)
	}
	op.argumentsJSON = append(op.argumentsJSON,
		fmt.Sprintf("\"%s_%s\":%s", jsonName, op.identifier, string(argJSON)))
	return nil
}

type iCanGenerateQuery interface {
	generateQuery(isTestMode ...bool) error
	getQuery() string
	getQueryParams() []string
	getQueryVariablesJSON() []string
}

func (cont *queryPart) generateSubQuery(cgq iCanGenerateQuery) (err error) {
	err = cgq.generateQuery()
	if err != nil {
		return newErr(ErrOperationCouldNotBeGenerated, err)
	}
	cont.queries = append(cont.queries, cgq.getQuery())
	cont.queryParams = append(cont.queryParams, cgq.getQueryParams()...)
	cont.queryVariablesJSON = append(cont.queryVariablesJSON, cgq.getQueryVariablesJSON()...)
	return nil
}

func (cont *queryPart) buildQuery(isMutation bool) {
	//build an overall query from the QueryOps in the container
	var csb strings.Builder

	switch isMutation {
	case true:
		csb.WriteString("mutation")
	case false:
		csb.WriteString("query")
	}

	if len(cont.queryParams) > 0 {
		csb.WriteString(fmt.Sprintf("(%s)", strings.Join(cont.queryParams, ",")))
	}

	csb.WriteString(fmt.Sprintf("{%s}", strings.Join(cont.queries, ",")))
	cont.query = csb.String()

	//reset the strings.Builder so we can use it again for the arguments
	csb.Reset()

	csb.WriteString(fmt.Sprintf("{%s}", strings.Join(cont.queryVariablesJSON, ",")))
	cont.argumentJSON = csb.String()
}

func lowerFirst(in string) string {
	sin := []byte(in)
	if len(in) > 1 {
		sin[0] = byte(unicode.ToLower(rune(sin[0])))
	}
	return string(sin)
}

var (
	ErrArgumentCouldNotBeEncodedAsJSON    = errors.New("argument could not be encoded as JSON")
	ErrArgumentError                      = errors.New("unkown argument error")
	ErrArgumentsCannotBeEmpty             = errors.New("arguments cannot be empty")
	ErrArgumentWasNotASliceOrArray        = errors.New("argument was not a slice or array")
	ErrArgumentWasNotAStruct              = errors.New("argument was not a struct")
	ErrJSONTagMissingOrNoValue            = errors.New("json tag was missing or did not have a value")
	ErrNoIsResultTagFoundInFieldsOfStruct = errors.New("no isResult tag found in fields of struct")
	ErrOperationCouldNotBeGenerated       = errors.New("operation could not be generated")
	ErrResultTypeNotHandled               = errors.New("result type not handled")
	ErrResultWasNotASliceOrArray          = errors.New("result was not a slice or array")
	ErrResultWasNotPointer                = errors.New("result was not pointer")
	ErrResultWasNotPointerToStruct        = errors.New("result was not pointer to struct")
)
