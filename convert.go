package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"strings"
)

// Convert takes a Go file and converts it to a Cap'n Proto schema
func Convert(goFile string) (*bytes.Buffer, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, goFile, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	conf := types.Config{Importer: importer.For("source", nil)}
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}
	_, err = conf.Check("", fset, []*ast.File{node}, info)
	if err != nil {
		return nil, err
	}

	structs := make(map[string][]string)
	interfaces := make(map[string][]string)
	importedStructs := make(map[string][]string)

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.GenDecl:
			for _, spec := range x.Specs {
				switch ts := spec.(type) {
				case *ast.TypeSpec:
					switch tp := ts.Type.(type) {
					case *ast.StructType:
						var fields []string
						for i, field := range tp.Fields.List {
							fieldType := typeToString(field.Type, info)
							fieldNames := make([]string, len(field.Names))
							for j, name := range field.Names {
								fieldNames[j] = toLowerCamelCase(name.Name)
							}
							if len(fieldNames) == 0 {
								// Embedded struct
								fieldNames = append(fieldNames, toLowerCamelCase(fieldType))
							}
							fields = append(fields, fmt.Sprintf("%s @%d :%s;", strings.Join(fieldNames, ", "), i, goToCapnp(fieldType)))
						}
						structs[ts.Name.Name] = fields
					case *ast.InterfaceType:
						var methods []string
						for i, method := range tp.Methods.List {
							switch ft := method.Type.(type) {
							case *ast.FuncType:
								var params, results []string
								for _, param := range ft.Params.List {
									paramType := typeToString(param.Type, info)
									paramNames := make([]string, len(param.Names))
									for j, name := range param.Names {
										paramNames[j] = toLowerCamelCase(name.Name)
									}
									params = append(params, fmt.Sprintf("%s :%s", strings.Join(paramNames, ", "), goToCapnp(paramType)))
								}
								for j, result := range ft.Results.List {
									resultType := typeToString(result.Type, info)
									results = append(results, fmt.Sprintf("result%d :%s", j, goToCapnp(resultType)))
								}
								methods = append(methods, fmt.Sprintf("%s @%d (%s) -> (%s);", toLowerCamelCase(method.Names[0].Name), i, strings.Join(params, ", "), strings.Join(results, ", ")))
							}
						}
						interfaces[ts.Name.Name] = methods
					}
				}
			}
		case *ast.SelectorExpr:
			// Handle imported structs
			typeName := x.Sel.Name
			obj := info.Uses[x.Sel]
			if obj != nil {
				structType := obj.Type().Underlying().(*types.Struct)
				var fields []string
				for i := 0; i < structType.NumFields(); i++ {
					field := structType.Field(i)
					fieldType := field.Type().String()
					fields = append(fields, fmt.Sprintf("%s @%d :%s;", toLowerCamelCase(field.Name()), i, goToCapnp(fieldType)))
				}
				importedStructs[typeName] = fields
			}
		}
		return true
	})

	var capnpSchema bytes.Buffer

	// Add necessary headers
	capnpSchema.WriteString("using Go = import \"/go.capnp\";\n")
	capnpSchema.WriteString(fmt.Sprintf("@0x%s\n", genNewCapnpId()))
	capnpSchema.WriteString("$Go.package(\"contract_impl\");\n")
	capnpSchema.WriteString("$Go.import(\"github.com/RedUndercover/capnconvert-go/testdata/contract\");\n\n")

	// Generate Cap'n Proto schema for structs
	for structName, fields := range structs {
		capnpSchema.WriteString(fmt.Sprintf("struct %s {\n", structName))
		for _, field := range fields {
			capnpSchema.WriteString(fmt.Sprintf("  %s\n", field))
		}
		capnpSchema.WriteString("}\n")
	}

	// Generate Cap'n Proto schema for imported structs
	for structName, fields := range importedStructs {
		capnpSchema.WriteString(fmt.Sprintf("struct %s {\n", structName))
		for _, field := range fields {
			capnpSchema.WriteString(fmt.Sprintf("  %s\n", field))
		}
		capnpSchema.WriteString("}\n")
	}

	// Generate Cap'n Proto schema for interfaces
	for interfaceName, methods := range interfaces {
		capnpSchema.WriteString(fmt.Sprintf("interface %s {\n", interfaceName))
		for _, method := range methods {
			capnpSchema.WriteString(fmt.Sprintf("  %s\n", method))
		}
		capnpSchema.WriteString("}\n")
	}

	return &capnpSchema, nil
}
