package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/machship-oss/go-gql"
)

func main() {
	// jsnTest()
	// jsnTest2()
	// gqltest()
	// queryTests()
	mutationTests()
}

func queryTests() {
	g := gql.NewGraphQL("http://10.200.100.69:8080/graphql", nil)
	// g := gql.NewGraphQL("http://localhost:1923", nil)

	qry := &gql.QueryContainer{
		QueryOps: []*gql.QueryOp{
			{
				ResultObject: &MonthMultiResult{},
				Arguments: &gql.QueryArguments{
					First: gql.NewInt(5),
					Order: &MonthOrder{
						Asc: OC_Month_Number, //todo: this doesn't appear to work. Fix it
					},
				},
			},
			{
				ResultObject: &SeasonMultiResult{},
				Arguments: &gql.QueryArguments{
					First: gql.NewInt(2),
				},
			},
		},
		GetOps: []*gql.GetOp{
			{
				ResultObject: &Month{},
				Arguments: &MonthGetArgs{
					ID: gql.NewID("0x220beaf"),
				},
			},
		},
		AggregateOps: []*gql.AggregateOp{
			{
				ResultObject: &MonthAggregateResult{},
				Arguments: &gql.AggregateArguments{
					Filter: &MonthFilter{
						Has: HC_Month_Season,
					},
				},
			},
		},
	}

	err := g.BulkQuery(context.Background(), qry)
	if err != nil {
		log.Fatal(err)
	}

	qryRes1 := qry.QueryOps[0].ResultObject.(*MonthMultiResult)
	qryRes2 := qry.QueryOps[1].ResultObject.(*SeasonMultiResult)
	getRes1 := qry.GetOps[0].ResultObject.(*Month)
	aggRes1 := qry.AggregateOps[0].ResultObject.(*MonthAggregateResult)

	fmt.Println("------------------------------------")
	fmt.Println("------------------------------------")
	fmt.Println()

	for _, x := range *qryRes1 {
		fmt.Printf("%+v\n", x.ID)
		fmt.Printf("\t%+v\n", x.Name)
	}

	fmt.Println()

	for _, x := range *qryRes2 {
		fmt.Printf("%+v\n", x.ID)
		fmt.Printf("\t%+v\n", x.Name)
	}

	fmt.Println()

	if getRes1 != nil {
		x := getRes1
		fmt.Printf("%+v\n", x.ID)
		fmt.Printf("\t%+v\n", x.Name)
	}

	fmt.Println()

	if aggRes1 != nil {
		x := aggRes1
		fmt.Printf("\t%+v\n", x.Count)
	}
}

func mutationTests() {
	g := gql.NewGraphQL("http://10.200.100.69:8080/graphql", nil)
	// g := gql.NewGraphQL("http://localhost:1923", nil)

	idSpring := "0x220bee3"
	idSummer := "0x220bee0"
	idNote := "0x220bef1"

	qry := &gql.MutationContainer{
		AddOps: []*gql.AddOp{
			{
				ResultObject: &NoteMutationResult{},
				Arguments: &AddNoteInputs{
					&AddNoteInput{
						NoteFields: NoteFields{
							Note:           gql.NewString("This is a note for Spring"),
							DateCreatedUTC: gql.NewTime(time.Now()),
						},
						BelongsTo: &gql.BaseRef{
							ID: gql.NewID(idSpring),
						},
					},
				},
			},
			{
				ResultObject: &NoteMutationResult{},
				Arguments: &AddNoteInputs{
					&AddNoteInput{
						NoteFields: NoteFields{
							Note:           gql.NewString("This is a different note for Summer"),
							DateCreatedUTC: gql.NewTime(time.Now()),
						},
						BelongsTo: &gql.BaseRef{
							ID: gql.NewID(idSummer),
						},
					},
				},
			},
		},
		UpdateOps: []*gql.UpdateOp{
			{
				ResultObject: &NoteMutationResult{},
				Arguments: &UpdateNoteInput{
					Filter: &NoteFilter{
						ID: gql.NewID(idNote),
					},
					Set: &NotePatch{
						NoteFields: NoteFields{
							Note: gql.NewString("updated note"),
						},
					},
				},
			},
		},
		DeleteOps: []*gql.DeleteOp{
			{
				ResultObject: &NoteDeleteResult{},
				Arguments: &NoteDeleteArguments{
					Filter: &NoteFilter{
						ID: gql.NewID(idNote),
					},
				},
			},
		},
	}

	err := g.BulkMutation(context.Background(), qry)
	if err != nil {
		log.Fatal(err)
	}

	addRes1 := qry.AddOps[0].ResultObject.(*NoteMutationResult)
	addRes2 := qry.AddOps[1].ResultObject.(*NoteMutationResult)
	udtRes1 := qry.UpdateOps[0].ResultObject.(*NoteMutationResult)
	delRes1 := qry.DeleteOps[0].ResultObject.(*NoteDeleteResult)

	fmt.Println("------------------------------------")
	fmt.Println("------------------------------------")
	fmt.Println()

	for _, x := range addRes1.Notes {
		fmt.Printf("%+v\n", x.ID)
		fmt.Printf("\t%+v\n", x.Note)
	}

	fmt.Println()

	for _, x := range addRes2.Notes {
		fmt.Printf("%+v\n", x.ID)
		fmt.Printf("\t%+v\n", x.Note)
	}

	fmt.Println()

	for _, x := range udtRes1.Notes {
		fmt.Printf("%+v\n", x.ID)
		fmt.Printf("\t%+v\n", x.Note)
	}

	fmt.Println()

	for _, x := range delRes1.Notes {
		fmt.Printf("%+v\n", x.ID)
		fmt.Printf("\t%+v\n", x.Note)
	}
}
