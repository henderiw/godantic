package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/henderiw/godantic/pkg/genvalidate/types"
)

const validationMarker = "// +generate:validate"

type StructInfo struct {
	Name   string
	Fields []FieldInfo
}

type EnumInfo struct {
	Name          string
	Type          string
	AllowedValues []string
}

type FieldInfo struct {
	Name            string
	Type            ast.Expr
	ValidationRules []types.ValidationRule
	Node            *ast.File
}

type FileInfo struct {
	Path    string
	Package string
	Structs []StructInfo
	Enums   []EnumInfo
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/main.go <path>")
		os.Exit(1)
	}
	path := os.Args[1]
	validategenerator := NewGenerator(path)
	validategenerator.Generate()

	/*
		x := networkv1alpha1.Dummy(1)
		if err := x.Validate(); err != nil {
			panic(err)
		}
	*/
}

func NewGenerator(path string) *Generator {
	return &Generator{
		path:     path,
		registry: types.InitValidationRuleRegistry(),
	}
}

type Generator struct {
	path     string
	registry map[string]types.ValidatorRuleParser
}

func (r *Generator) Generate() {
	filepath.Walk(r.path, func(path string, info os.FileInfo, err error) error {
		if err != nil || !strings.HasSuffix(path, ".go") {
			return nil
		}
		fileInfo, err := r.processFile(path)
		if err != nil {
			fmt.Println("Error processing", path, ":", err)
		}
		r.generateValidationCode(fileInfo)
		return nil
	})

}

func (r *Generator) processFile(path string) (*FileInfo, error) {
	fs := token.NewFileSet()
	node, err := parser.ParseFile(fs, path, nil, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return nil, err
	}
	fileInfo := &FileInfo{
		Path:    path,
		Package: node.Name.Name, // Extract package name
		Structs: []StructInfo{},
		Enums:   []EnumInfo{},
	}

	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			if genDecl.Doc == nil {
				continue
			}
			// initialize generateValidation
			generateValidation := false
			// Check if marker comment is present
			for _, comment := range genDecl.Doc.List {
				if strings.TrimSpace(comment.Text) == validationMarker {
					generateValidation = true
					break
				}
			}

			if !generateValidation {
				continue
			}

			switch typeDecl := typeSpec.Type.(type) {
			case *ast.StructType:
				// Handle Structs
				var fields []FieldInfo
				for _, field := range typeDecl.Fields.List {
					if len(field.Names) == 0 {
						continue
					}

					var validationRules []types.ValidationRule
					skip := false
					if field.Doc != nil {
						for _, comment := range field.Doc.List {
							if strings.HasPrefix(comment.Text, "// +validate(") {
								funcName, attr, err := r.parseValidation(comment.Text)
								if err != nil {
									panic(err)
								}

								if funcName == "skip" {
									skip = true
									break
								}

								validationRule, err := r.parseValidationRule(funcName, attr)
								if err != nil {
									panic(err)
								}
								validationRules = append(validationRules, validationRule)
							}
						}
					}

					if !skip {
						fields = append(fields, FieldInfo{
							Name:            field.Names[0].Name,
							Type:            field.Type,
							ValidationRules: validationRules,
							Node:            node,
						})
					}
				}
				fileInfo.Structs = append(fileInfo.Structs, StructInfo{Name: typeSpec.Name.Name, Fields: fields})
			case *ast.Ident:
				// Handle Enum-like Types (Alias of string, int, etc.)
				if baseType, isAlias := detectTypeAlias(node, typeSpec.Name.Name); isAlias {
					allowedValues := extractEnumValues(node, typeSpec.Name.Name)
					fileInfo.Enums = append(fileInfo.Enums, EnumInfo{
						Name:          typeSpec.Name.Name,
						Type:          baseType,
						AllowedValues: allowedValues,
					})
				}
			}
		}
	}

	return fileInfo, nil
}

func (r *Generator) parseValidation(comment string) (string, string, error) {
	comment = strings.TrimSpace(comment)
	if !strings.HasPrefix(comment, "// +validate(") {
		return "", "", fmt.Errorf("unexpected comment, expected `// +validate` prefix got: %s", comment)
	}

	// Extract content inside `()`
	start := strings.Index(comment, "(")
	end := strings.LastIndex(comment, ")")
	if start == -1 || end == -1 || start > end {
		fmt.Println("Invalid validate syntax")
		return "", "", fmt.Errorf("invalid validate syntax got: %s", comment)
	}

	content := comment[start+1 : end] // Get `length(min = 4, max = 5, equal = 5)`

	// Split at the first `(` to get `length`
	parts := strings.SplitN(content, "(", 2)
	funcName := strings.TrimSpace(parts[0])
	if len(parts) == 1 {
		return funcName, "", nil
	}
	attrs := strings.TrimSuffix(parts[1], ")")
	return funcName, attrs, nil
}

func (r *Generator) parseValidationRule(funcName, attrs string) (types.ValidationRule, error) {
	// Check if the function is registered
	if parserFunc, exists := r.registry[funcName]; exists {
		return parserFunc(attrs)
	}

	return nil, fmt.Errorf("unsupported validator: %s", funcName)
}

func exprToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.ArrayType:
		return "[]" + exprToString(t.Elt)
	case *ast.MapType:
		return "map[" + exprToString(t.Key) + "]" + exprToString(t.Value)
	default:
		return "unknown"
	}
}

func (r *Generator) generateValidationCode(fileInfo *FileInfo) {
	outputFile := strings.TrimSuffix(fileInfo.Path, ".go") + "_validate.go"
	var sb strings.Builder
	sb.WriteString("// GENERATED CODE - DO NOT EDIT\n")
	sb.WriteString(fmt.Sprintf("package %s\n\n", fileInfo.Package)) // Use actual package name
	sb.WriteString("import (\n")
	sb.WriteString("\t\"fmt\"\n")
	//sb.WriteString("\t\"errors\"\n")
	sb.WriteString("\t)\n")

	for _, enum := range fileInfo.Enums {
		//var sb strings.Builder
		sb.WriteString(fmt.Sprintf("func (r %s) Validate() error {\n", enum.Name))
		sb.WriteString(generateEnumValidation(enum.Name, enum.Type, enum.AllowedValues))
		sb.WriteString("}\n")
		//fmt.Println(sb.String())
	}
	/*
		for _, s := range fileInfo.Structs {
			var sb strings.Builder
			//sb.WriteString(fmt.Sprintf("func (s *%s) Validate() error {\n", s.Name))
			//sb.WriteString("\tvar errs error\n")

			for _, f := range s.Fields {
				//for _, validationRule := range f.ValidationRules {
				//	fmt.Printf("\tfield: %s type: %s rules: %v\n", f.Name, f.Type, validationRule.String())
				//}
				rulePresent := false
				for _, rule := range f.ValidationRules {
					rulePresent = true
					sb.WriteString(rule.String())
				}
				if !rulePresent {
					nested := generateNestedStructs(f)
					if len(nested) > 0 {
						sb.WriteString(nested)
					}
				}


				//	if f.Validate != "" {
				//		fmt.Printf("\tif err := validate%s(s.%s); err != nil {\n", strings.Title(f.Validate), f.Name)
				//		fmt.Println("\t\terrors = append(errors, err)")
				//		fmt.Println("\t}")
				//	}

			}

			sb.WriteString("\tif errs != nil{ return errs }\n")
			sb.WriteString("\treturn nil\n")
			sb.WriteString("}\n")
			fmt.Println(sb.String())
		}
	*/
	if len(fileInfo.Enums) > 0 {
		err := os.WriteFile(outputFile, []byte(sb.String()), 0644)
		if err != nil {
			fmt.Println("Error writing validation file:", err)
			return
		}

		// Run gofmt
		formatGoFile(outputFile)

		fmt.Println("Generated validation file:", outputFile)
	}

}

func generateNestedStructs(fieldInfo FieldInfo) string {
	typ := fieldInfo.Type
	pointer := false
	var innerType ast.Expr = typ
	if starExpr, ok := typ.(*ast.StarExpr); ok {
		pointer = true
		innerType = starExpr.X
	}

	var sb strings.Builder

	switch t := innerType.(type) {
	case *ast.StructType, *ast.SelectorExpr:
		sb.WriteString(fmt.Sprintf("\tSTRUCT %s\n", fieldInfo.Name))
		// EXPAND ENUM HERE

		// If it's a struct, call Validate() on it
		if pointer {
			sb.WriteString(fmt.Sprintf("\tif %s != nil { if err := %s.Validate(); err != nil { errs = errors.Join(errs, err) } }\n",
				fieldInfo.Name, fieldInfo.Name))
		} else {
			sb.WriteString(fmt.Sprintf("\tif err := %s.Validate(); err != nil { errs = errors.Join(errs, err }\n", fieldInfo.Name))
		}

	case *ast.ArrayType:
		sb.WriteString(fmt.Sprintf("\tARRAY %s\n", fieldInfo.Name))
		// If it's a slice, iterate and call Validate()
		elementType := t.Elt
		pointerElement := false

		if starExpr, ok := elementType.(*ast.StarExpr); ok {
			pointerElement = true
			elementType = starExpr.X
		}

		if _, isStruct := elementType.(*ast.StructType); isStruct {
			// EXPAND ENUM HERE
			if pointerElement {
				sb.WriteString(fmt.Sprintf("\tfor _, item := range %s { if item != nil { if err := item.Validate(); err != nil { errs = errors.Join(errs, err } } }\n",
					fieldInfo.Name))
			} else {
				sb.WriteString(fmt.Sprintf("\tfor _, item := range %s { if err := item.Validate(); err != nil { errs = errors.Join(errs, err } }\n",
					fieldInfo.Name))
			}
		}

	case *ast.MapType:
		sb.WriteString(fmt.Sprintf("\tMAP %s\n", fieldInfo.Name))
		// If it's a map, iterate over values and call Validate()
		valueType := t.Value
		pointerValue := false

		if starExpr, ok := valueType.(*ast.StarExpr); ok {
			pointerValue = true
			valueType = starExpr.X
		}

		if _, isStruct := valueType.(*ast.StructType); isStruct {
			// EXPAND ENUM HERE
			if pointerValue {
				sb.WriteString(fmt.Sprintf("\tfor _, value := range %s { if value != nil { if err := value.Validate(); err != nil { errs = errors.Join(errs, err } } }\n",
					fieldInfo.Name))
			} else {
				sb.WriteString(fmt.Sprintf("\tfor _, value := range %s { if err := value.Validate(); err != nil { errs = errors.Join(errs, err } }\n",
					fieldInfo.Name))
			}
		}
	}
	return sb.String()

}

// detectTypeAlias checks if a type is an alias of string.
func detectTypeAlias(node ast.Node, typeName string) (string, bool) {
	found := false
	baseType := ""

	ast.Inspect(node, func(n ast.Node) bool {
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok || typeSpec.Name.Name != typeName {
			return true
		}

		// Check if it's a base type
		if ident, ok := typeSpec.Type.(*ast.Ident); ok {
			baseType = ident.Name
			found = true
		}

		return false
	})
	return baseType, found
}

// generateEnumValidation generates an enum validation function
func generateEnumValidation(fieldName, typeName string, allowedValues []string) string {
	if len(allowedValues) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\tvalid := map[%s]struct{}{", typeName))
	for _, v := range allowedValues {
		sb.WriteString(fmt.Sprintf("\t%s: {}, ", v))
	}
	sb.WriteString("}\n")

	// Generate validation check
	v := "s"
	if strings.HasPrefix(typeName, "int") || strings.HasPrefix(typeName, "uint") {
		v = "d"
	}
	sb.WriteString(fmt.Sprintf(
		`if _, ok := valid[%s(r)]; !ok {
    		return fmt.Errorf("invalid value for %s: %%%s", r)
		}
		return nil`,
		typeName, fieldName, v))

	return sb.String()
}

// extractEnumValues extracts allowed values for an enum-like type (e.g., AdminState).

func extractEnumValues(node *ast.File, typeName string) []string {
	var values []string
	var lastKnownType string // Store the last explicitly declared type
	currentIotaValue := 0    // Track the `iota` counter
	iotaDeclared := false    // Detect if an `iota` block is ongoing

	ast.Inspect(node, func(n ast.Node) bool {
		decl, ok := n.(*ast.GenDecl)
		if !ok || decl.Tok != token.CONST {
			return true
		}

		for _, spec := range decl.Specs {
			valSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			// Check if this constant has an explicit type
			if valSpec.Type != nil {
				if ident, ok := valSpec.Type.(*ast.Ident); ok {
					lastKnownType = ident.Name // Store the type
				}
			}

			// If the type is missing, use the last known type
			if lastKnownType == typeName {
				for i := range valSpec.Names {
					if i < len(valSpec.Values) {
						switch v := valSpec.Values[i].(type) {
						case *ast.BasicLit: // Directly assigned integer values
							values = append(values, v.Value)
							iotaDeclared = false // Stop tracking `iota`

						case *ast.Ident: // Handles explicit `iota`
							if v.Name == "iota" {
								values = append(values, fmt.Sprintf("%d", currentIotaValue))
								currentIotaValue++ // Increment for next constant
								iotaDeclared = true
							}

						default:
							iotaDeclared = false
						}
					} else if iotaDeclared { // If `iota` is inherited, increment manually
						values = append(values, fmt.Sprintf("%d", currentIotaValue))
						currentIotaValue++
					}
				}
			}
		}

		return true
	})

	return values
}

// Run gofmt on the generated file
func formatGoFile(filename string) {
	cmd := exec.Command("gofmt", "-w", filename)
	if err := cmd.Run(); err != nil {
		fmt.Println("Error formatting file with gofmt:", err)
	}
}
