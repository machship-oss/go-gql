package gql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context/ctxhttp"
)

// Client is a GraphQL client.
type GraphQL struct {
	url        string // GraphQL server URL.
	httpClient *http.Client
}

// NewClient creates a GraphQL client targeting the specified GraphQL server URL.
// If httpClient is nil, then http.DefaultClient is used.
func NewGraphQL(url string, httpClient *http.Client) *GraphQL {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &GraphQL{
		url:        url,
		httpClient: httpClient,
	}
}

func (c *GraphQL) Query(ctx context.Context, op *QueryOp) (err error) {
	cont := &QueryContainer{
		QueryOps: []*QueryOp{
			op,
		},
	}

	return c.BulkQuery(ctx, cont)
}

func (c *GraphQL) Get(ctx context.Context, op *GetOp) (err error) {
	cont := &QueryContainer{
		GetOps: []*GetOp{
			op,
		},
	}

	return c.BulkQuery(ctx, cont)
}

func (c *GraphQL) Aggregate(ctx context.Context, op *AggregateOp) (err error) {
	cont := &QueryContainer{
		AggregateOps: []*AggregateOp{
			op,
		},
	}

	return c.BulkQuery(ctx, cont)
}

func (c *GraphQL) Add(ctx context.Context, op *AddOp) (err error) {
	cont := &MutationContainer{
		AddOps: []*AddOp{
			op,
		},
	}

	return c.BulkMutation(ctx, cont)
}

func (c *GraphQL) Update(ctx context.Context, op *UpdateOp) (err error) {
	cont := &MutationContainer{
		UpdateOps: []*UpdateOp{
			op,
		},
	}

	return c.BulkMutation(ctx, cont)
}

func (c *GraphQL) Delete(ctx context.Context, op *DeleteOp) (err error) {
	cont := &MutationContainer{
		DeleteOps: []*DeleteOp{
			op,
		},
	}

	return c.BulkMutation(ctx, cont)
}

func (c *GraphQL) BulkQuery(ctx context.Context, cont *QueryContainer) (err error) {
	err = cont.generateQuery()
	if err != nil {
		return err
	}

	//need to 'do' the query, then match up the results from the data set
	outTypes := make(map[string]interface{})
	for _, op := range cont.QueryOps {
		outTypes[op.identifier] = op.ResultObject
	}
	for _, op := range cont.GetOps {
		outTypes[op.identifier] = op.ResultObject
	}
	for _, op := range cont.AggregateOps {
		outTypes[op.identifier] = op.ResultObject
	}

	err = c.do(ctx, cont.queryPart, outTypes)
	if err != nil {
		return err
	}

	return nil
}

func (c *GraphQL) BulkMutation(ctx context.Context, cont *MutationContainer) (err error) {
	err = cont.generateQuery()
	if err != nil {
		return err
	}

	//need to 'do' the query, then match up the results from the data set
	outTypes := make(map[string]interface{})
	for _, op := range cont.AddOps {
		outTypes[op.identifier] = op.ResultObject
	}
	for _, op := range cont.UpdateOps {
		outTypes[op.identifier] = op.ResultObject
	}
	for _, op := range cont.DeleteOps {
		outTypes[op.identifier] = op.ResultObject
	}

	err = c.do(ctx, cont.queryPart, outTypes)
	if err != nil {
		return err
	}

	return nil
}

func (c *GraphQL) DoRaw(ctx context.Context, query string, vars map[string]interface{}, outputTypes map[string]interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()

	in := struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables,omitempty"`
	}{
		Query:     query,
		Variables: vars,
	}

	fmt.Println(in.Query)
	varJsonBytes, _ := json.Marshal(&in.Variables)
	fmt.Println()
	fmt.Println(string(varJsonBytes))
	fmt.Println()

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(in)
	if err != nil {
		return err
	}

	resp, err := ctxhttp.Post(ctx, c.httpClient, c.url, "application/json", &buf)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	var bodyBuf bytes.Buffer
	tee := io.TeeReader(resp.Body, &bodyBuf)
	bodyBytes, _ := ioutil.ReadAll(tee)
	fmt.Println(string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(&bodyBuf)
		return fmt.Errorf("non-200 OK status code: %v body: %q", resp.Status, body)
	}

	var out struct {
		Data   *json.RawMessage `json:"data"`
		Errors GraphqlErrors    `json:"errors"`
		// Extensions interface{} // todo: add this later. should conform to some standard
	}

	err = json.NewDecoder(&bodyBuf).Decode(&out)
	if err != nil {
		// TODO: Consider including response body in returned error, if deemed helpful.
		return err
	}

	var outData interface{}

	if out.Data != nil {
		// err := jsonutil.UnmarshalGraphQL(*out.Data, output) //TBC on how to do this with dynamic results
		err := json.Unmarshal(*out.Data, &outData) //TBC on how to do this with dynamic results
		if err != nil {
			// TODO: Consider including response body in returned error, if deemed helpful.
			return err
		}
	}

	if len(out.Errors) > 0 {
		return out.Errors
	}

	//this may panic, so there is a defer recover on this method
	msgMap := outData.(map[string]interface{})
	for k, mm := range msgMap {
		//find it in the output map:
		ot, ok := outputTypes[k]
		if ok {
			bts, _ := json.Marshal(mm)
			_ = json.Unmarshal(bts, &ot)
		}
	}

	return nil
}

func (c *GraphQL) do(ctx context.Context, cont queryPart, outputTypes map[string]interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()

	var vars map[string]interface{}
	err = json.Unmarshal([]byte(cont.argumentJSON), &vars)
	if err != nil {
		return fmt.Errorf("couldn't convert variables to map")
	}

	in := struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables,omitempty"`
	}{
		Query:     cont.query,
		Variables: vars,
	}

	fmt.Println(in.Query)
	varJsonBytes, _ := json.Marshal(&in.Variables)
	fmt.Println()
	fmt.Println(string(varJsonBytes))
	fmt.Println()

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(in)
	if err != nil {
		return err
	}

	resp, err := ctxhttp.Post(ctx, c.httpClient, c.url, "application/json", &buf)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	var bodyBuf bytes.Buffer
	tee := io.TeeReader(resp.Body, &bodyBuf)
	bodyBytes, _ := ioutil.ReadAll(tee)
	fmt.Println(string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(&bodyBuf)
		return fmt.Errorf("non-200 OK status code: %v body: %q", resp.Status, body)
	}

	var out struct {
		Data   *json.RawMessage `json:"data"`
		Errors GraphqlErrors    `json:"errors"`
		// Extensions interface{} // todo: add this later. should conform to some standard
	}

	err = json.NewDecoder(&bodyBuf).Decode(&out)
	if err != nil {
		// TODO: Consider including response body in returned error, if deemed helpful.
		return err
	}

	var outData interface{}

	if out.Data != nil {
		// err := jsonutil.UnmarshalGraphQL(*out.Data, output) //TBC on how to do this with dynamic results
		err := json.Unmarshal(*out.Data, &outData) //TBC on how to do this with dynamic results
		if err != nil {
			// TODO: Consider including response body in returned error, if deemed helpful.
			return err
		}
	}

	if len(out.Errors) > 0 {
		return out.Errors
	}

	//this may panic, so there is a defer recover on this method
	msgMap := outData.(map[string]interface{})
	for k, mm := range msgMap {
		//find it in the output map:
		ot, ok := outputTypes[k]
		if ok {
			bts, _ := json.Marshal(mm)
			_ = json.Unmarshal(bts, &ot)
		}
	}

	return nil
}

// errors represents the "errors" array in a response from a GraphQL server.
// If returned via error interface, the slice is expected to contain at least 1 element.
//
// Specification: https://facebook.github.io/graphql/#sec-Errors.
type GraphqlErrors []struct {
	Message   string `json:"message"`
	Locations []struct {
		Line   int `json:"line,omitempty"`
		Column int `json:"column,omitempty"`
	} `json:"locations,omitempty"`
}

// Error implements error interface.
func (e GraphqlErrors) Error() string {
	jsn, _ := json.Marshal(e)
	return string(jsn)
}
