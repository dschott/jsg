package main

import (
	"reflect"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

type File struct {
	Version string
	Pkg     string
	Types   []*Type
}

type Type struct {
	Schema jsonschema.Schema
	Name   string
	Pkg    string
	Kind   reflect.Kind
	Elem   *Type
	Fields []Field
}

type Field struct {
	Name string
	Type Type
	Tag  reflect.StructTag
}
