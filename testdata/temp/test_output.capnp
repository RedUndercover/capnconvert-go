using Go = import "/go.capnp";
@0xf1190b5d41ed11ef
$Go.package("contract_impl");
$Go.import("github.com/RedUndercover/capnconvert-go/testdata/contract");

struct SimpleStruct {
  id @0 :Int32;
  name @1 :Text;
}
struct NestedStruct {
  simple @0 :SimpleStruct;
  value @1 :Float64;
}
struct ComplexStruct {
  embeddedImport @0 :ImportedStruct;
  anotherField @1 :Text;
  nested @2 :AnotherImportedStruct;
}
struct ImportedStruct {
  field1 @0 :Text;
  field2 @1 :Int32;
}
struct AnotherImportedStruct {
  nestedField @0 :ImportedStruct;
  value @1 :Float64;
}
interface MyInterface {
  foo @0 (x :Int32) -> (result0 :SimpleStruct);
  bar @1 (y :Text) -> (result0 :Text);
  baz @2 (arr :List(Int32)) -> (result0 :List(Text));
  multipleReturns @3 () -> (result0 :Int32, result1 :Text);
  nestedReturn @4 () -> (result0 :NestedStruct);
  importTest @5 () -> (result0 :ImportedStruct);
}
interface AnotherInterface {
  complexMethod @0 (x :ImportedStruct) -> (result0 :AnotherImportedStruct, result1 :Text);
}
