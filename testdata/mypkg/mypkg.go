package mypkg

type ImportedStruct struct {
	Field1 string
	Field2 int
}

type AnotherImportedStruct struct {
	NestedField ImportedStruct
	Value       float64
}
