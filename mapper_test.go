package main

import (
	"encoding/json"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var schemas = mustLoadSchemas(`{
    "$id": "foo",
    "type": "object",
    "required": [ "name" ],
    "properties": {
        "name": { "type": "string" },
        "bars": {
            "type": "array",
            "items": { "$ref": "bar" }
        }
    }
}`, `{
    "$id": "bar",
    "type": "object",
    "required": [ "name" ],
    "properties": {
        "name": { "type": "string" }
    }
}`)

func Test_Mapper(t *testing.T) {
	var mapper Mapper
	typ, err := mapper.Map(schemas["foo"])
	require.NoError(t, err)

	assert.Equal(t, "main", typ.Pkg)
	assert.Equal(t, "Foo", typ.Name)
}

func Test_Mapper_ToIdentifier(t *testing.T) {
	check := func(s string, want string, initialisms []string) {
		t.Helper()
		var m Mapper
		for _, ini := range initialisms {
			m.AddInitialism(ini)
		}
		got := m.ToIdentifier(s)
		assert.Equal(t, want, got, "input: "+s)
	}

	check("a", "A", nil)
	check("A", "A", nil)
	check("a1 ", "A1", nil)
	check("1a", "_1A", nil)
	check("_a", "A", nil)
	check("_A", "A", nil)
	check("_1", "1", nil)
	check("_1a", "1A", nil)
	check("_1A", "1A", nil)
	check(".1a", "1A", nil)
	check("._1a", "1A", nil)
	check(" _1a ", "1A", nil)
	check("_a ", "A", nil)
	check("aa", "Aa", nil)
	check("aaa", "Aaa", nil)
	check("aA ", "AA", nil)
	check("_abc123_xyz ", "Abc123Xyz", nil)
	check("abc", "ABC", []string{"ABC"})
	check("123abc", "_123ABC", []string{"ABC"})
	check("abc123", "ABC123", []string{"ABC"})
	check("abc123", "Abc123", []string{"B"})
	check("aBc123", "ABC123", []string{"BC"})
}

func mustLoadSchemas(jsonDocs ...string) map[string]*jsonschema.Schema {
	compiler := jsonschema.NewCompiler()
	var ids []string
	for _, jsonDoc := range jsonDocs {
		var doc map[string]any
		if err := json.Unmarshal([]byte(jsonDoc), &doc); err != nil {
			panic(err)
		}
		id := doc["$id"].(string)
		ids = append(ids, id)
		if err := compiler.AddResource(id, doc); err != nil {
			panic(err)
		}
	}
	schemas := make(map[string]*jsonschema.Schema)
	for _, id := range ids {
		schema, err := compiler.Compile(id)
		if err != nil {
			panic(err)
		}
		schemas[id] = schema
	}
	return schemas
}
