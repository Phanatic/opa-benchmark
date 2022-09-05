package lib

import (
	"testing"
	"encoding/json"
	"github.com/open-policy-agent/opa/rego"
)

func BenchmarkOPA_PrepareAndEvaluate(b *testing.B) {
	for i:=0; i < b.N; i++ {
		Run()
	}
}
func BenchmarkOPA_EvaluateAllowAll_UsingPreparedQuery(b *testing.B) {
	query, _ :=  PrepareAllowAllQuery()
	for i:=0; i < b.N; i++ {
		var input map[string]interface{}
		inputJson := `{
    "user": "sam",
    "action": "read",
    "table": "customers",
    "columns": [
        "id",
        "name"
    ]
}`
		json.Unmarshal([]byte(inputJson), &input)
		Evaluate(query, input)
	}
}

func BenchmarkOPA_EvaluateUsingPreparedQuery(b *testing.B) {
	query, _ :=  PrepareQuery()
	for i:=0; i < b.N; i++ {
		var input map[string]interface{}
		inputJson := `{
    "user": "sam",
    "action": "read",
    "table": "customers",
    "columns": [
        "id",
        "name"
    ]
}`
		json.Unmarshal([]byte(inputJson), &input)
		Evaluate(query, input)
	}
}

func BenchmarkOPA_Evaluate_PreparedQuery_PreparedInput_NoColumns(b *testing.B) {
	query, _ :=  PrepareQuery()
	var input map[string]interface{}
	inputJson := `{
    "user": "sam",
    "action": "read",
    "table": "customers",
    "columns": []
}`
	json.Unmarshal([]byte(inputJson), &input)
	for i:=0; i < b.N; i++ {
		Evaluate(query, input)
	}
}

func BenchmarkOPA_Evaluate_PreparedQuery_PreparedInput_SingleColumn(b *testing.B) {
	query, _ :=  PrepareQuery()

	var input map[string]interface{}
	inputJson := `{
    "user": "sam",
    "action": "read",
    "table": "customers",
    "columns": [
        "id"
    ]
}`
	json.Unmarshal([]byte(inputJson), &input)
	for i:=0; i < b.N; i++ {
		Evaluate(query, input)
	}
}

func BenchmarkOPA_Evaluate_PreparedQuery_PreparedInput_MultipleColumns(b *testing.B) {
	query, _ :=  PrepareQuery()

	var input map[string]interface{}
	inputJson := `{
    "user": "sam",
    "action": "read",
    "table": "customers",
    "columns": [
        "id",
		"name"
    ]
}`
	json.Unmarshal([]byte(inputJson), &input)
	for i:=0; i < b.N; i++ {
		Evaluate(query, input)
	}
}

func BenchmarkOPA_Evaluate_PreparedQuery_EvaledInput_MultipleColumns(b *testing.B) {
	query, _ :=  PrepareQuery()

	var input map[string]interface{}
	inputJson := `{
    "user": "sam",
    "action": "read",
    "table": "customers",
    "columns": [
        "id",
		"name"
    ]
}`
	json.Unmarshal([]byte(inputJson), &input)
	opt := rego.EvalInput(input)
	for i:=0; i < b.N; i++ {
		EvaluateWithOption(query, opt)
	}
}