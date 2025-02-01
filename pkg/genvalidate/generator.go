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
	Name            string
	Fields          []FieldInfo
	HasNestedStruct bool
	HasValidationRules bool
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
	NestedStruct    bool
}

type FileInfo struct {
	Path             string
	Package          string
	Structs          []StructInfo
	Enums            []EnumInfo
	HasNestedStructs bool
	HasValidationRules bool
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/main.go <path>")
		os.Exit(1)
	}
	path := os.Args[1]
	validategenerator := NewGenerator(path)
	validategenerator.Generate()

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

	fileHasNestedStructs := false
	fileHasValidationRules := false
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
				var hasNestedStruct bool
				var hasValidationRules bool
				var fields []FieldInfo
				for _, field := range typeDecl.Fields.List {
					if len(field.Names) == 0 {
						continue
					}

					var validationRules []types.ValidationRule
					var nestedStruct bool
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
								hasValidationRules = true
								fileHasValidationRules = true
							} else {
								nestedStruct = isNestedStructOrEnum(field.Type, node)
								if nestedStruct {
									hasNestedStruct = true
									fileHasNestedStructs = true
								}
							}
						}
					}

					if !skip {
						fields = append(fields, FieldInfo{
							Name:            field.Names[0].Name,
							Type:            field.Type,
							ValidationRules: validationRules,
							Node:            node,
							NestedStruct:    nestedStruct,
						})
					}
				}
				fileInfo.Structs = append(fileInfo.Structs, StructInfo{
					Name:            typeSpec.Name.Name,
					Fields:          fields,
					HasNestedStruct: hasNestedStruct,
					HasValidationRules: hasValidationRules,
				})
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
	fileInfo.HasNestedStructs = fileHasNestedStructs
	fileInfo.HasValidationRules = fileHasValidationRules
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

func (r *Generator) generateValidationCode(fileInfo *FileInfo) {
	outputFile := strings.TrimSuffix(fileInfo.Path, ".go") + "_validate.go"
	var sb strings.Builder
	sb.WriteString("// GENERATED CODE - DO NOT EDIT\n")
	sb.WriteString(fmt.Sprintf("package %s\n\n", fileInfo.Package)) // Use actual package name
	sb.WriteString("import (\n")
	if len(fileInfo.Enums) > 0 || fileInfo.HasValidationRules {
		sb.WriteString("\t\"fmt\"\n")
	}
	if fileInfo.HasNestedStructs {
		fmt.Println(fileInfo.Path, len(fileInfo.Structs), fileInfo.HasNestedStructs)
		sb.WriteString("\t\"errors\"\n")
	}
	sb.WriteString("\t)\n")

	for _, enumInfo := range fileInfo.Enums {
		sb.WriteString(fmt.Sprintf("func (r %s) Validate() error {\n", enumInfo.Name))
		sb.WriteString(generateEnumValidation(enumInfo.Name, enumInfo.Type, enumInfo.AllowedValues))
		sb.WriteString("}\n")
	}
	for _, schemaInfo := range fileInfo.Structs {
		sb.WriteString(fmt.Sprintf("func (r *%s) Validate() error {\n", schemaInfo.Name))
		if schemaInfo.Name == "ConditionedStatus" {
			fmt.Println("ConditionedStatus", schemaInfo)
		}
		if schemaInfo.HasNestedStruct {
			sb.WriteString("\tvar errs error\n")
		}

		for _, fieldInfo := range schemaInfo.Fields {
			for _, rule := range fieldInfo.ValidationRules {
				fieldName := fieldInfo.Name
				fieldNameCode := fmt.Sprintf("r.%s", fieldName)
				if isPointerType(fieldInfo.Type) {
					sb.WriteString(fmt.Sprintf("if r.%s != nil {\n", fieldName))
					fieldNameCode = fmt.Sprintf("*r.%s", fieldName)
					fieldName = "*" + fieldName // Dereference pointer for validation
				}
				
				sb.WriteString(rule.ExpandCode(fieldName, fieldNameCode)) // Expand the code based on the rule
				if isPointerType(fieldInfo.Type) {
					sb.WriteString("}\n") // Close the pointer check block
				}
			}
			// nested code generation is implicitly enabled
			// when a struct exists we generate the nested validation rules
			if fieldInfo.NestedStruct {
				sb.WriteString(generateNestedStructs(fieldInfo, fieldInfo.Type, fmt.Sprintf("r.%s", fieldInfo.Name)))
			}
		}
		if schemaInfo.HasNestedStruct {
			sb.WriteString("\tif errs != nil{ return errs }\n")
		}
		sb.WriteString("\treturn nil\n")
		sb.WriteString("}\n")
	}
	if len(fileInfo.Enums) > 0 || len(fileInfo.Structs) > 0 {
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

func isNestedStructOrEnum(expr ast.Expr, node *ast.File) bool {
	switch t := expr.(type) {
	case *ast.StarExpr: // Pointer to a struct (e.g., *MyStruct)
		return isNestedStructOrEnum(t.X, node)
	case *ast.StructType: // Direct struct type
		return true
	case *ast.SelectorExpr: // External package struct (e.g., mypkg.MyStruct)
		return true
	case *ast.Ident:
		return isDeclaredAsStructOrEnum(t.Name, node)
	case *ast.ArrayType: // Array (e.g., []MyStruct)
		fmt.Println("array type")
		return isNestedStructOrEnum(t.Elt, node) // Recursively check element type
	case *ast.MapType: // Map (e.g., map[string]MyStruct)
		return isNestedStructOrEnum(t.Value, node) // Recursively check map value type
	default:
		return false
	}
}

func generateNestedStructs(fieldInfo FieldInfo, expr ast.Expr, fieldName string) string {
	var sb strings.Builder

	switch t := expr.(type) {
	case *ast.StarExpr:
		// If it's a pointer, wrap validation inside `if != nil`
		sb.WriteString(fmt.Sprintf("if %s != nil {\n", fieldName))
		sb.WriteString(generateNestedStructs(fieldInfo, t.X, fieldName))
		sb.WriteString("}\n")

	// the ast.ident we blindly use since we have done the validation before (ast.Ident is s struct in the same file)
	case *ast.StructType, *ast.SelectorExpr, *ast.Ident:
		// If it's a struct, call Validate()
		sb.WriteString(fmt.Sprintf("if err := %s.Validate(); err != nil {\n", fieldName))
		sb.WriteString("\terrs = errors.Join(errs, err)\n")
		sb.WriteString("}\n")

	case *ast.ArrayType:
		// If it's an array/slice, iterate and call Validate()
		iteratorVar := "item"
		sb.WriteString(fmt.Sprintf("for _, %s := range %s {\n", iteratorVar, fieldName))
		sb.WriteString(generateNestedStructs(fieldInfo, t.Elt, iteratorVar))
		sb.WriteString("}\n")

	case *ast.MapType:
		// If it's a map, iterate over values and call Validate()
		iteratorVar := "value"
		sb.WriteString(fmt.Sprintf("for _, %s := range %s {\n", iteratorVar, fieldName))
		sb.WriteString(generateNestedStructs(fieldInfo, t.Value, iteratorVar))
		sb.WriteString("}\n")
	}

	return sb.String()
}

// isIdentDeclaredAsStruct check if an `ast.Ident` is a Struct
func isDeclaredAsStructOrEnum(typeName string, node *ast.File) bool {
	for _, decl := range node.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if typeSpec.Name.Name == typeName {
						_, isStruct := typeSpec.Type.(*ast.StructType)
						_, isEnum := detectTypeAlias(node, typeName)
						return isStruct || isEnum
					}
				}
			}
		}
	}
	return false
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


func isPointerType(expr ast.Expr) bool {
	_, ok := expr.(*ast.StarExpr)
	return ok
}