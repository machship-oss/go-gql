package main

import (
	"fmt"
	"strings"
)

func (o *Object) generateFileBytes(pkg string) []byte {
	sb := &strings.Builder{}

	sb.WriteString(fmt.Sprintf(`
	package %s

	`, pkg))

	fmtString := ""
	if o.IsEnumType {
		fmtString = `"fmt"`
	}

	gqlString := ""
	// for _, _ = range o.Fields {
	// if f.FieldType != FT_Object {
	gqlString = `"github.com/machship-oss/go-gql"`
	// 	break
	// }
	// }

	if fmtString != "" || gqlString != "" {
		sb.WriteString(fmt.Sprintf(`
	import (
		%s

		%s
	)

	`, fmtString, gqlString))
	}

	o.oStruct(sb)
	o.oFields(sb)
	o.oGetArgs(sb)
	o.oRelationships(sb)
	o.oInputUpdate(sb)
	o.oPatch(sb)
	o.oInputAdd(sb)
	o.oAggregate(sb)
	o.oMutationResult(sb)
	o.oFilter(sb)
	o.oOrderable(sb)
	o.oOrder(sb)
	o.oHasFilter(sb)
	o.oPasswordCheck(sb)
	o.oWebhook(sb)
	// o.oDgraphMethods(sb)

	return []byte(sb.String())
}

func (o *Object) getBaseString(str string) string {
	if o.ImplementsBase {
		return str
	}
	return ""
}

func (o *Object) oStruct(sb *strings.Builder) {
	var str string
	if o.Description != "" {
		str = fmt.Sprintf(`
	/*
	%s
	*/`, o.Description)

		sb.WriteString(str)
	}

	str = fmt.Sprintf(`
	type %[1]s struct {
		%[2]s
		*%[1]sFields

		*%[1]sRelationships
	}

	func (ar *%[1]s) SingleResultName() string {
		return "%[1]s"
	}

	type %[1]sMultiResult []*%[1]s

	func (ar *%[1]sMultiResult) MultiResultName() string {
		return "%[1]s"
	}
	`, o.Name, o.getBaseString("gql.Base"))

	sb.WriteString(str)
}

func (o *Object) oFields(sb *strings.Builder) {
	str := fmt.Sprintf(`
	type %sFields struct {
	`, o.Name)
	sb.WriteString(str)

	for _, f := range o.Fields {
		if f.IsSecret {
			continue //we do not add secrets as queryable
		}

		if f.InverseObject == nil {
			if f.Description != "" {
				str = fmt.Sprintf(`
			// %s
			`, f.Description)

				sb.WriteString(str)
			}

			sb.WriteString(
				fmt.Sprintf("	%s	*%s	`json:\"%s,omitempty\"`\n",
					f.Name, string(f.FieldType), lowerFirst(f.Name)))
		}
	}

	str = `
	}

	`
	sb.WriteString(str)
}

func (o *Object) oGetArgs(sb *strings.Builder) {
	str := fmt.Sprintf(`
	type %[1]sGetArgs struct {
		%[2]s
	`, o.Name, o.getBaseString("gql.BaseGet"))
	sb.WriteString(str)

	for _, f := range o.Fields {
		if f.IsSecret {
			continue //we do not add secrets as queryable
		}

		if f.InverseObject == nil {
			if f.Description != "" {
				str = fmt.Sprintf(`
			// %s
			`, f.Description)

				sb.WriteString(str)
			}

			sb.WriteString(
				fmt.Sprintf("	%s	*%s	`json:\"%s,omitempty\"`\n",
					f.Name, string(f.FieldType), lowerFirst(f.Name)))
		}
	}

	str = fmt.Sprintf(`
	}


	func (t *%[1]sGetArgs) IsArg() {}

	func (t *%[1]sGetArgs) IsGet() {}
	`, o.Name)
	sb.WriteString(str)
}

func (o *Object) oRelationships(sb *strings.Builder) {
	str := fmt.Sprintf(`
	type %sRelationships struct {
	`, o.Name)
	sb.WriteString(str)

	for _, f := range o.Fields {
		if f.IsSecret {
			continue //we do not add secrets as queryable
		}

		if f.InverseObject != nil {
			if f.Description != "" {
				str = fmt.Sprintf(`
			// %s
			`, f.Description)

				sb.WriteString(str)
			}

			if !f.IsPlural {
				sb.WriteString(
					fmt.Sprintf("	%s	*%s	`hasInverse:\"%s\" json:\"%s,omitempty\"`\n",
						f.Name, f.InverseTypeName, f.InverseName, lowerFirst(f.Name)))
			}
			if f.IsPlural {
				sb.WriteString(
					fmt.Sprintf("	%s	[]*%s	`hasInverse:\"%s\" json:\"%s,omitempty\"`\n",
						f.Name, f.InverseTypeName, f.InverseName, lowerFirst(f.Name)))

			}
		}

	}

	str = `
	}
	`
	sb.WriteString(str)
}

func (o *Object) oInputUpdate(sb *strings.Builder) {
	str := fmt.Sprintf(`
	type Update%sInput struct {
	`, o.Name)
	sb.WriteString(str)

	str =
		fmt.Sprintf(`
		Filter *%[1]sFilter |json:"filter"|
		Set *%[1]sPatch |json:"set,omitempty"|
		Remove *%[1]sPatch |json:"remove,omitempty"|
	}

	func (t *Update%[1]sInput) IsArg()    {}
	
	func (t *Update%[1]sInput) IsUpdate() {}
	`,
			o.Name)

	str = strings.ReplaceAll(str, "|", "`")

	sb.WriteString(str)
}

func (o *Object) oPatch(sb *strings.Builder) {
	str := fmt.Sprintf(`
	type %[1]sPatch struct {
		%[2]s
		%[1]sFields
	`, o.Name, o.getBaseString("gql.BasePatch"))
	sb.WriteString(str)

	for _, f := range o.Fields {
		if f.InverseObject != nil {
			if !f.IsPlural {
				sb.WriteString(
					fmt.Sprintf("	%s	*gql.BaseRef	`json:\"%s,omitempty\"`\n",
						f.Name, lowerFirst(f.Name)))
			}
			if f.IsPlural {
				sb.WriteString(
					fmt.Sprintf("	%s	[]*gql.BaseRef	`json:\"%s,omitempty\"`\n",
						f.Name, lowerFirst(f.Name)))

			}
		}
		if f.IsSecret {
			sb.WriteString(
				fmt.Sprintf("	%s	*gql.String	`json:\"%s,omitempty\"`\n",
					f.Name, lowerFirst(f.Name)))
		}
	}

	str = `
	}
	`
	sb.WriteString(str)
}

func (o *Object) oInputAdd(sb *strings.Builder) {
	str := fmt.Sprintf(`
	type Add%[1]sInput struct {
		%[2]s
		%[1]sFields
	`, o.Name, o.getBaseString("gql.AddBaseInput"))
	sb.WriteString(str)

	for _, f := range o.Fields {
		if f.InverseObject != nil {
			if !f.IsPlural {
				sb.WriteString(
					fmt.Sprintf("	%s	*gql.BaseRef	`json:\"%s,omitempty\"`\n",
						f.Name, lowerFirst(f.Name)))
			}
			if f.IsPlural {
				sb.WriteString(
					fmt.Sprintf("	%s	[]*gql.BaseRef	`json:\"%s,omitempty\"`\n",
						f.Name, lowerFirst(f.Name)))

			}
		}
		if f.IsSecret {
			sb.WriteString(
				fmt.Sprintf("	%s	*gql.String	`json:\"%s,omitempty\"`\n",
					f.Name, lowerFirst(f.Name)))
		}
	}

	str = fmt.Sprintf(`
	}

	type Add%[1]sInputs []*Add%[1]sInput

	func (t *Add%[1]sInputs) IsArg() {}
	
	func (t *Add%[1]sInputs) AddName() string {
		return "Add%[1]sInput"
	}
	`, o.Name)
	sb.WriteString(str)
}

func (o *Object) oAggregate(sb *strings.Builder) {
	str := fmt.Sprintf(`
	type %[1]sAggregateResult struct {
		%[3]s
	}

	func (ar *%[1]sAggregateResult) AggregateResultName() string {
		return "%[2]s"
	}
	`,
		o.Name, upperFirst(o.Name), "gql.BaseAggregate")

	sb.WriteString(str)
}

func (o *Object) oMutationResult(sb *strings.Builder) {
	str := fmt.Sprintf(`
	type %[1]sMutationResult struct {
		%[1]ss []*%[1]s |isResult:"true" json:"%[2]s,omitempty"|
		NumberAffected gql.Int              |json:"numUids,omitempty"|
	}
	
	type %[1]sDeleteResult struct {
		%[1]sMutationResult
		Message gql.String |json:"msg,omitempty"|
	}
	
	func (r %[1]sMutationResult) MutationName() string {
		return "%[3]s"
	}

	type %[1]sDeleteArguments struct {
		Filter *%[1]sFilter |json:"filter"|
	}
	
	func (t *%[1]sDeleteArguments) IsArg() {}
	
	func (t *%[1]sDeleteArguments) DeleteName() string {
		return "%[1]sFilter"
	}

	func (t *%[1]sDeleteArguments) DeleteFilter() gql.IFilter {
		return t.Filter
	}

	`,
		o.Name, lowerFirst(o.Name), upperFirst(o.Name))

	str = strings.ReplaceAll(str, "|", "`")

	sb.WriteString(str)
}

func (o *Object) oFilter(sb *strings.Builder) {
	str := fmt.Sprintf(`
	type %[1]sFilter struct {
		%[2]s
	`, o.Name, o.getBaseString("gql.BaseFilter"))
	sb.WriteString(str)

	for _, f := range o.Fields {
		if f.IsSecret {
			continue //we do not add secrets as queryable
		}

		if !(strings.ToUpper(f.Name) == "ID" || f.IsSearchable) {
			continue //we only add the ID field or searchable fields
		}

		if f.InverseObject == nil {
			sb.WriteString(
				fmt.Sprintf("	%s	*%s	`json:\"%s,omitempty\"`\n",
					f.Name, string(f.FieldType), lowerFirst(f.Name)))
		}
	}

	str =
		fmt.Sprintf(`
		Has %[1]sHasChoice |json:"has,omitempty"|
		And []*%[1]sFilter |json:"and,omitempty"|
		Or []*%[1]sFilter |json:"or,omitempty"|
		Not []*%[1]sFilter |json:"not,omitempty"|
	}

	func (t *%[1]sFilter) GetName() string {
		return "%[1]sFilter"
	}

	func (t *%[1]sFilter) IsFilter() {}
	
	func (t *%[1]sFilter) IsArg() {}

`,
			o.Name)

	str = strings.ReplaceAll(str, "|", "`")
	sb.WriteString(str)
}

func (o *Object) oOrderable(sb *strings.Builder) {
	// fnd := false
	// for _, f := range o.Fields {
	// 	if f.IsSecret {
	// 		continue //we do not add secrets as queryable
	// 	}

	// 	// if f.IsSearchable {
	// 	fnd = true
	// 	// }
	// }
	// if !fnd {
	// 	return
	// }

	typeName := "string"
	if o.ImplementsBase {
		typeName = "gql.BaseOrderChoice"
	}

	str := fmt.Sprintf(`
	type %[1]sOrderChoice %[2]s

	const (
	`, o.Name, typeName)
	sb.WriteString(str)

	for _, f := range o.Fields {
		if f.IsSecret {
			continue //we do not add secrets as queryable
		}

		if f.IsSearchable {
			sb.WriteString(
				fmt.Sprintf("	OC_%s_%s %sOrderChoice = \"%s\"\n",
					o.Name, f.Name, o.Name, lowerFirst(f.Name)))
		}
	}

	str = `
	)
	`
	sb.WriteString(str)
}

func (o *Object) oOrder(sb *strings.Builder) {
	str := fmt.Sprintf(`
	type %sOrder struct {
	`, o.Name)
	sb.WriteString(str)

	str =
		fmt.Sprintf(`
		Asc %[1]sOrderChoice |json:"asc,omitempty"|
		Desc %[1]sOrderChoice |json:"desc,omitempty"|
		Then *%[1]sOrder |json:"then,omitempty"|
	}

	func (t *%[1]sOrder) GetName() string {
		return "%[1]sOrder"
	}
	
	func (t *%[1]sOrder) IsOrder() {}
	
`,
			o.Name)

	str = strings.ReplaceAll(str, "|", "`")
	sb.WriteString(str)
}

func (o *Object) oHasFilter(sb *strings.Builder) {
	// if len(o.Fields) == 0 {
	// 	return
	// }

	typeName := "string"
	if o.ImplementsBase {
		typeName = "gql.BaseHasChoice"
	}

	str := fmt.Sprintf(`
	type %[1]sHasChoice %[2]s

	const (
	`, o.Name, typeName)
	sb.WriteString(str)

	for _, f := range o.Fields {
		if f.IsSecret {
			continue //we do not add secrets as queryable
		}

		sb.WriteString(
			fmt.Sprintf("	HC_%s_%s %sHasChoice = \"%s\"\n",
				o.Name, f.Name, o.Name, lowerFirst(f.Name)))
	}

	str = `
	)
	`
	sb.WriteString(str)
}

func (o *Object) oPasswordCheck(sb *strings.Builder) {
	for _, f := range o.Fields {
		if f.IsSecret {
			str := fmt.Sprintf(`
			func (o *%s) CheckPassword(%s *gql.String) (isCorrect bool) {
				//do some password stuff
				return false
			}
			`,
				o.Name, o.Name)

			sb.WriteString(str)

			return //we only allow a single 'secret' field, as does Dgraph
		}
	}
}

func (o *Object) oWebhook(sb *strings.Builder) {
	if o.Webhook != nil {
		str := fmt.Sprintf(`
		func (o *%s) GetWebhookType() string {
			return "%s"
		}

		func (o *%s) GetWebhookSource() string {
			return "%s"
		}

		func (o *%s) GetObject() interface{} {
			return o
		}
		`,
			o.Name, o.Webhook.Type, o.Name, o.Webhook.Source, o.Name)

		sb.WriteString(str)
	}
}

//////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////

func (o *Object) oDgraphMethods(sb *strings.Builder) {
	str := fmt.Sprintf(`
		func (db *GraphQL) Query%ss(filter string) ([]*%s, error) {
			var q struct {
				data []*%s |graphql:"query%s(input: $input)"|
			}
		
			err := db.Query(&q, nil)
			if err != nil {
				return nil, err
			}
		
			return q.data, nil
		}
		
		func (db *GraphQL) Get%sByID(id string) (*%s, error) {
			var q struct {
				data *%s
			}
		
			err := db.Query(&q, db.GetFilterByID(id))
			if err != nil {
				return nil, err
			}
		
			return q.data, nil
		}
		
		func (db *GraphQL) Add%s(input *Add%sInput) (*%s, error) {
			var m struct {
				addedData *%s |graphql:"add%s(input: $input)"|
			}
		
			err := db.Add(&m, input)
			if err != nil {
				return nil, err
			}
		
			return m.addedData, nil
		}
		
		func (db *GraphQL) Update%s(input *Update%sInput) (*%s, error) {
			var m struct {
				updatedData *%s |graphql:"update%s(input: $input)"|
			}
		
			err := db.Update(&m, input)
			if err != nil {
				return nil, err
			}
		
			return m.updatedData, nil
		}
		
		func (db *GraphQL) Delete%s(filter *%sFilter) ([]*%s, error) {
			var m struct {
				deletedData []*%s |graphql:"delete%s(input: $input)"|
			}
		
			err := db.Delete(&m, filter)
			if err != nil {
				return nil, err
			}
		
			return m.deletedData, nil
		}
		`,
		o.Name, o.Name, o.Name, o.Name, o.Name,
		o.Name, o.Name, o.Name, o.Name, o.Name,
		o.Name, o.Name, o.Name, o.Name, o.Name,
		o.Name, o.Name, o.Name, o.Name, o.Name,
		o.Name, o.Name)

	str = strings.ReplaceAll(str, "|", "`")

	sb.WriteString(str)
}

func generateBase(pkg string) []byte {
	var sb strings.Builder

	def := fmt.Sprintf(`
package %s

import (
	"github.com/machship-oss/go-gql"
)

type GraphQL struct {
	*gql.GraphQL
}


`, pkg)

	def = strings.ReplaceAll(def, "|", "`")

	sb.WriteString(def)

	return []byte(sb.String())
}
