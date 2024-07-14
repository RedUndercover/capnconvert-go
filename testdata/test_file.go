package main

import "github.com/RedUndercover/capnconvert-go/testdata/mypkg"

type SimpleStruct struct {
	ID   int
	Name string
}

type NestedStruct struct {
	Simple SimpleStruct
	Value  float64
}

type MyInterface interface {
	Foo(x int) SimpleStruct
	Bar(y string) string
	Baz(arr []int) []string
	MultipleReturns() (int, error)
	NestedReturn() NestedStruct
	ImportTest() mypkg.ImportedStruct
}

type ComplexStruct struct {
	EmbeddedImport mypkg.ImportedStruct
	AnotherField   string
	Nested         mypkg.AnotherImportedStruct
}

type AnotherInterface interface {
	ComplexMethod(x mypkg.ImportedStruct) (mypkg.AnotherImportedStruct, error)
}
