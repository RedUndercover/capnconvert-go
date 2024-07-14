struct ComplexStruct {
  EmbeddedImport @0 :ImportedStruct;
  AnotherField @1 :Text;
  Nested @2 :AnotherImportedStruct;
}
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
interface MyInterface {
  Foo @0 (x :Int32) -> (result0 :SimpleStruct);
  Bar @1 (y :Text) -> (result0 :Text);
  Baz @2 (arr :List(Int32)) -> (result0 :List(Text));
  MultipleReturns @3 () -> (result0 :Int32, result1 :Text);
  NestedReturn @4 () -> (result0 :NestedStruct);
  ImportTest @5 () -> (result0 :ImportedStruct);
}
interface AnotherInterface {
  ComplexMethod @0 (x :ImportedStruct) -> (result0 :AnotherImportedStruct, result1 :Text);
}
