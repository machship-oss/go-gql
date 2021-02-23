package main

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/dgraph-io/gqlparser/v2/ast"
	"github.com/dgraph-io/gqlparser/v2/parser"
	"github.com/dgraph-io/gqlparser/v2/validator"
	"github.com/pkg/errors"
)

func main() {
	st := time.Now()
	do()
	fmt.Println("Finished in:", time.Since(st))
}

func do() {
	fileName := "./newschema.graphql"
	byts, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	schema := string(byts)

	doc, gqlErr := parser.ParseSchemas(validator.Prelude, &ast.Source{Input: schema})
	if gqlErr != nil {
		log.Fatal(errors.Wrap(gqlErr, "while parsing GraphQL schema"))
	}

	ol := build(&doc.Definitions)

	ol.backfill()

	pkgname := "main"
	dirname := "./output/" + time.Now().Format("20060102150405") + "/" + pkgname

	_ = os.MkdirAll(dirname, 0777)

	baseStr := generateBase(pkgname)
	// fmt.Println(string(baseStr))
	bbyts, err := format.Source(baseStr)
	if err != nil {
		log.Fatal("berr:\n", err)
	}
	_ = ioutil.WriteFile(fmt.Sprintf("%s/base.go", dirname), bbyts, 0777)

	for _, o := range ol.Objects {
		if o.IsIBaseEnum {
			continue
		}

		objStr := o.generateFileBytes(pkgname)
		// fmt.Println(string(objStr))
		obyts, err := format.Source(objStr)
		if err != nil {
			fmt.Println(string(objStr))
			log.Fatal("oerr:\n", err)
		}
		_ = ioutil.WriteFile(fmt.Sprintf("%s/%s.go", dirname, lowerFirst(o.Name)), obyts, 0777)
	}
}

func upperFirst(in string) string {
	if strings.ToLower(in) == "id" {
		return "ID"
	}

	sin := []byte(in)
	if len(in) > 1 {
		sin[0] = byte(unicode.ToUpper(rune(sin[0])))
	}
	return string(sin)
}

func lowerFirst(in string) string {
	if strings.ToLower(in) == "id" {
		return "id"
	}

	sin := []byte(in)
	if len(in) > 1 {
		sin[0] = byte(unicode.ToLower(rune(sin[0])))
	}
	return string(sin)
}

type ObjectList struct {
	Objects map[string]*Object
}

type Object struct {
	Name           string
	Description    string
	IsIBaseEnum    bool
	IsEnumType     bool
	IsInterface    bool
	ImplementsBase bool
	Interfaces     []Interface
	Fields         map[string]*Field
	Webhook        *Webhook
}

type Webhook struct {
	Type   string
	Source string
}

type Interface struct {
	Name   string
	Object *Object
}

type Field struct {
	Name              string
	Description       string
	IsBaseObjectField bool
	IsSecret          bool
	IsRequired        bool
	IsPlural          bool
	FieldType         FieldType
	InverseName       string
	InverseTypeName   string
	InverseObject     *Object
	IsSearchable      bool
}

type FieldType string

const (
	FT_ID      FieldType = "gql.ID"
	FT_Bool    FieldType = "gql.Bool"
	FT_Float64 FieldType = "gql.Float64"
	FT_Int64   FieldType = "gql.Int64"
	FT_Int     FieldType = "gql.Int64"
	FT_String  FieldType = "gql.String"
	FT_Time    FieldType = "gql.Time"
	FT_Point   FieldType = "gql.Point"
	FT_Object  FieldType = "OBJECT"
)

type BOF_BaseObjectField string

const (
	BOF_ID          BOF_BaseObjectField = "ID"
	BOF_DateCreated BOF_BaseObjectField = "DateCreatedUTC"
)

func checkIfFieldBaseObjectField(in string) bool {
	return false
	//todo: ignore the rest for now, will put the fields on the objects
	switch strings.ToUpper(in) {
	case strings.ToUpper(string(BOF_ID)), strings.ToUpper(string(BOF_DateCreated)):
		return true
	}
	return false
}

func (ol *ObjectList) backfill() {
	for _, o := range ol.Objects {
		for fk, f := range o.Fields {
			if f.IsBaseObjectField {
				delete(o.Fields, fk)
				continue
			}
			if f.InverseTypeName != "" {
				//link the inverse object
				if inv, ok := ol.Objects[strings.ToUpper(f.InverseTypeName)]; ok {
					f.InverseObject = inv
				} else {
					fmt.Fprintf(os.Stderr, "object '%s' wanted inverse object '%s' for field '%s' but it couldn't be found\n", o.Name, f.InverseName, f.Name)
				}
			}
		}
	}

	//now check that all inverse relationships have their inverse defined
	for _, o := range ol.Objects {
		for _, f := range o.Fields {
			if f.FieldType == FT_Object && f.InverseName == "" {
				fmt.Fprintf(os.Stderr, "object '%s' has field '%s' which should have an inverse field defined\n", o.Name, f.Name)
			}
			if invObj := f.InverseObject; invObj != nil {
				//get the inverse object, then get the labelled inverse field and check that it is this object
				if invFld, fok := invObj.Fields[strings.ToUpper(f.InverseName)]; fok {
					if invFld.InverseObject != o {
						fmt.Fprintf(os.Stderr, "object '%s' has field '%s' which should have an inverse field of type '%s' which is instead of type '%s'\n", invObj.Name, invFld.Name, o.Name, invFld.InverseName)
						fmt.Fprintf(os.Stderr, "object '%s' has field '%s' which should have an inverse field of type '%s' which is instead of type '%s'\n", invObj.Name, invFld.Name, o.Name, invFld.InverseName)
					}
				} else {
					fmt.Fprintf(os.Stderr, "object '%s' has field '%s' which should have an inverse field '%s' which could not be found on object '%s'\n", o.Name, f.Name, f.InverseName, invObj.Name)
				}
			}
		}
	}

	for _, o := range ol.Objects {
		for _, iface := range o.Interfaces {
			if iface.Name == "IBase" {
				//todo: ignore this for now, will actually put all the fields directly in the objects
				// o.ImplementsBase = true
			}

			//get the interface
			if ifaceObj, ok := ol.Objects[strings.ToUpper(iface.Name)]; ok {
				//add the interface's fields to the object
				for fk, fv := range ifaceObj.Fields {
					o.Fields[strings.ToUpper(fk)] = fv
				}
			} else {
				fmt.Fprintf(os.Stderr, "object '%s' wanted interface '%s' but it couldn't be found\n", o.Name, iface.Name)
			}
		}
	}
}

func build(defs *ast.DefinitionList) *ObjectList {
	ol := &ObjectList{
		Objects: make(map[string]*Object),
	}

	for _, x := range *defs {
		isInterface := false
		{
			//ignore:
			switch x.Kind {
			case "INTERFACE":
				isInterface = true
			case "SCALAR":
				continue
			}

			if strings.HasPrefix(x.Name, "__") {
				continue
			}
		}

		o := &Object{
			Name:        x.Name,
			IsInterface: isInterface,
			Description: x.Description,
			Fields:      make(map[string]*Field),
		}
		ol.Objects[strings.ToUpper(o.Name)] = o

		for _, iface := range x.Interfaces {
			o.Interfaces = append(o.Interfaces, Interface{
				Name: iface,
			})
		}

		if strings.HasSuffix(o.Name, "Enum") {
			o.IsEnumType = true
		}

		for _, y := range x.Fields {
			f := buildField(y)
			o.Fields[strings.ToUpper(f.Name)] = f
		}

		for _, y := range x.Directives {
			//this is a 'password' field
			if y.Name == "secret" {
				f := &Field{
					Name:       upperFirst(y.Arguments[0].Value.Raw),
					IsRequired: true,
					IsSecret:   true,
					FieldType:  FT_String,
				}
				o.Fields[strings.ToUpper(f.Name)] = f
			}

			//this is used internally to put properties and methods on the object which will be used in webhooks
			if y.Name == "webhook" {
				var src, typ string
				for _, a := range y.Arguments {
					switch a.Name {
					case "source":
						src = a.Value.Raw
					case "type":
						typ = a.Value.Raw
					}
				}

				if src != "" && typ != "" {
					o.Webhook = &Webhook{
						Source: src,
						Type:   typ,
					}
				}
			}
		}
	}

	//run through the new objects and make sure that all interfaces have the ID field
	for _, o := range ol.Objects {
		if !o.IsInterface {
			continue
		}

		fndID := false
		for _, f := range o.Fields {
			if strings.ToUpper(f.Name) == "ID" {
				fndID = true
			}
		}
		if !fndID {
			fmt.Fprintf(os.Stderr, "interface '%s' does not have an ID field which is not advisable\n", o.Name)
		}
	}

	return ol
}

func buildField(fd *ast.FieldDefinition) (f *Field) {
	f = &Field{
		Name:        upperFirst(fd.Name),
		IsRequired:  fd.Type.NonNull,
		Description: fd.Description,
	}

	if fd.Type.NamedType != "" {
		f.FieldType = getFieldType(fd.Type.NamedType)
	} else {
		f.FieldType = getFieldType(fd.Type.Elem.NamedType)
	}
	f.IsBaseObjectField = checkIfFieldBaseObjectField(f.Name)

	if f.FieldType == FT_Object {
		if fd.Type.NamedType != "" {
			f.InverseTypeName = (fd.Type.NamedType)
		} else {
			f.InverseTypeName = (fd.Type.Elem.NamedType)
		}
	}

	if fn := fd.Type.String(); strings.HasPrefix(fn, "[") && (strings.HasSuffix(fn, "]") || strings.HasSuffix(fn, "]!")) {
		f.IsPlural = true
	}

	for _, z := range fd.Directives {
		switch z.Name {
		case "search":
			f.IsSearchable = true
		case "hasInverse":
			for _, zz := range z.Arguments {
				switch zz.Name {
				case "field":
					f.InverseName = (zz.Value.Raw)
				}
			}
		}
	}

	return f
}

func getFieldType(in string) FieldType {
	stripped := strings.ReplaceAll(in, "[", "")
	stripped = strings.ReplaceAll(stripped, "]", "")
	stripped = strings.ReplaceAll(stripped, "!", "")

	switch stripped {
	case "ID":
		return FT_ID
	case "Boolean":
		return FT_Bool
	case "Float":
		return FT_Float64
	case "Int64":
		return FT_Int64
	case "Int":
		return FT_Int
	case "String":
		return FT_String
	case "DateTime":
		return FT_Time
	case "Point":
		return FT_Point
	}
	return FT_Object
}
