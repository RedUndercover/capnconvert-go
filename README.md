# Go to Cap'n Proto Converter

This project provides a CLI tool to convert Go structs and interfaces to Cap'n Proto schema definitions. It supports handling nested structs and interfaces, including those imported from other packages.

## Building the CLI

To build the CLI tool, run the following command from the project root:

```bash
go build -o convert ./cli.go
```

This will create an executable file named `convert`.

## Usage

### Basic Usage

To convert a Go file to a Cap'n Proto schema and print the output to the console, run:

```bash
./go_to_capnp path/to/your/file.go
```

### Example

Given a Go file `test_input.go`:

```go
package main

import "example.com/myproject/testdata/mypkg"

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
    EmbeddedImport mypkg.ImportedStruct
}

type ComplexStruct struct {
    EmbeddedImport mypkg.ImportedStruct
    AnotherField   string
    Nested         mypkg.AnotherImportedStruct
}

type AnotherInterface interface {
    ComplexMethod(x mypkg.ImportedStruct) (mypkg.AnotherImportedStruct, error)
}
```

Running the command:

```bash
./go_to_capnp test_input.go -o output.capnp
```

Will produce the following Cap'n Proto schema:

```plaintext
struct SimpleStruct {
  ID @0 :Int32;
  Name @1 :Text;
}
struct NestedStruct {
  Simple @0 :SimpleStruct;
  Value @1 :Float64;
}
struct ImportedStruct {
  Field1 @0 :Text;
  Field2 @1 :Int32;
}
struct AnotherImportedStruct {
  NestedField @0 :ImportedStruct;
  Value @1 :Float64;
}
struct ComplexStruct {
  EmbeddedImport @0 :ImportedStruct;
  AnotherField @1 :Text;
  Nested @2 :AnotherImportedStruct;
}
interface MyInterface {
  Foo @0 (x :Int32) -> (result0 :SimpleStruct);
  Bar @1 (y :Text) -> (result0 :Text);
  Baz @2 (arr :List(Int32)) -> (result0 :List(Text));
  MultipleReturns @3 () -> (result0 :Int32, result1 :Text);
  NestedReturn @4 () -> (result0 :NestedStruct);
  ImportTest @5 () -> (result0 :ImportedStruct);
  EmbeddedImport @6 () -> (result0 :ImportedStruct);
}
interface AnotherInterface {
  ComplexMethod @0 (x :ImportedStruct) -> (result0 :AnotherImportedStruct, result1 :Text);
}
```
