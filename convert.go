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

var goToCapnpType = map[string]string{
	"int":     "Int32",
	"int32":   "Int32",
	"int64":   "Int64",
	"float32": "Float32",
	"float64": "Float64",
	"string":  "Text",
	"bool":    "Bool",
	"error":   "Text",
}

func goToCapnp(goType string) string {
	if capnpType, found := goToCapnpType[goType]; found {
		return capnpType
	}
	if strings.HasPrefix(goType, "[]") {
		return fmt.Sprintf("List(%s)", goToCapnp(goType[2:]))
	}
	// if we have a . in the type, it's an imported type
	// it should be flattened to the last part of the type name in the schema
	if goTypeParts := strings.Split(goType, "."); len(goTypeParts) > 1 {
		//return the last part of the type
		return goTypeParts[len(goTypeParts)-1]
	}

	return goType // Assuming it's a custom type defined elsewhere
}

func typeToString(expr ast.Expr, info *types.Info) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.ArrayType:
		return "[]" + typeToString(t.Elt, info)
	case *ast.StarExpr:
		return "*" + typeToString(t.X, info)
	case *ast.SelectorExpr:
		// Handle imported types
		typeName := t.Sel.Name
		return typeName
	default:
		return fmt.Sprintf("%T", t)
	}
}

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
								fieldNames[j] = name.Name
							}
							if len(fieldNames) == 0 {
								// Embedded struct
								fieldNames = append(fieldNames, fieldType)
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
										paramNames[j] = name.Name
									}
									params = append(params, fmt.Sprintf("%s :%s", strings.Join(paramNames, ", "), goToCapnp(paramType)))
								}
								for j, result := range ft.Results.List {
									resultType := typeToString(result.Type, info)
									results = append(results, fmt.Sprintf("result%d :%s", j, goToCapnp(resultType)))
								}
								methods = append(methods, fmt.Sprintf("%s @%d (%s) -> (%s);", method.Names[0].Name, i, strings.Join(params, ", "), strings.Join(results, ", ")))
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
					fields = append(fields, fmt.Sprintf("%s @%d :%s;", field.Name(), i, goToCapnp(fieldType)))
				}
				importedStructs[typeName] = fields
			}
		}
		return true
	})

	var capnpSchema bytes.Buffer

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
