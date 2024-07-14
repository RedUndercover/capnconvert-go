// Description: This file contains functions that transform strings in some way
// The functions in this file are used to convert Go types to Cap'n Proto types
// and to convert Go struct field names to Cap'n Proto field names.
// The functions in this file are used in the convert.go file.
package main

import (
	"encoding/hex"
	"fmt"
	"go/ast"
	"go/types"
	"strings"
	"unicode"

	"github.com/google/uuid"
)

// Generates a new Cap'n Proto ID
// This is used to generate unique IDs for files, as capnp requires that all files have unique IDs
func genNewCapnpId() string {
	// Generate a new UUID
	u, _ := uuid.NewUUID() // if this errors out we're all screwed anyway

	// Get the first 8 bytes (64 bits) of the UUID
	first64Bits := u[:8]

	// Convert the first 8 bytes to a hex string
	hexString := hex.EncodeToString(first64Bits)

	return hexString
}

// Converts a string to lower camel case
// This is used to convert Go struct field names to Cap'n Proto field names
func toLowerCamelCase(s string) string {
	if s == "" {
		return s
	}
	if strings.ToUpper(s) == s {
		return strings.ToLower(s)
	}
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// Converts an AST type to a string representation
// This is used to convert Go types to Cap'n Proto types
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

// Map of Go types to Cap'n Proto types
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

// Converts a Go type to a Cap'n Proto type
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
