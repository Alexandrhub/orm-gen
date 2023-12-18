package genstorage

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

type FieldInfo struct {
	Name string
	Type string
}

type MethodInfo struct {
	Receiver string
	Name     string
}

type ReflectData struct {
	StructName string
	HasDBTag   bool
	Fields     []FieldInfo
	Methods    []MethodInfo
}

// ReflectFile функция извлечения информации о файле
func ReflectFile(fileName string) ([]ReflectData, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var reflectDataList []ReflectData

	ast.Inspect(
		node, func(n ast.Node) bool {
			switch t := n.(type) {
			case *ast.TypeSpec:
				var rData ReflectData
				rData.StructName = t.Name.Name

				s, ok := t.Type.(*ast.StructType)
				if ok {
					for _, field := range s.Fields.List {
						if len(field.Names) == 0 {
							continue
						}
						fieldName := field.Names[0].Name
						fieldType := fmt.Sprintf("%s", field.Type)
						if field.Tag != nil {
							tagValue := field.Tag.Value
							if strings.Contains(tagValue, "db:\"") {
								rData.HasDBTag = true
							}
						}
						rData.Fields = append(
							rData.Fields, FieldInfo{
								Name: fieldName,
								Type: fieldType,
							},
						)
					}
				}

				reflectDataList = append(reflectDataList, rData)

			case *ast.FuncDecl:
				if t.Recv != nil && len(t.Recv.List) == 1 {
					starExpr, ok := t.Recv.List[0].Type.(*ast.StarExpr)
					if ok {
						ident, ok := starExpr.X.(*ast.Ident)
						if ok {
							// Update the corresponding ReflectData for this receiver
							for idx, rData := range reflectDataList {
								if rData.StructName == ident.Name {
									reflectDataList[idx].Methods = append(
										reflectDataList[idx].Methods, MethodInfo{
											Receiver: ident.Name,
											Name:     t.Name.Name,
										},
									)
								}
							}
						}
					}
				}
			}
			return true
		},
	)

	return reflectDataList, nil
}

// GenerateMethods функция генерации методов
func GenerateMethods(data ReflectData) string {
	receiver := strings.ToLower(data.StructName[:1])
	builder := &strings.Builder{}

	if !methodExists(data.Methods, data.StructName, "TableName") {
		// Генерация метода TableName
		fmt.Fprintf(builder, "func (%s *%s) TableName() string {\n", receiver, data.StructName)
		fmt.Fprintf(builder, "\treturn \"%s\"\n", data.StructName)
		builder.WriteString("}\n\n")
	}

	if !methodExists(data.Methods, data.StructName, "OnCreate") {
		// Генерация метода OnCreate
		fmt.Fprintf(builder, "func (%s *%s) OnCreate() []string {\n", receiver, data.StructName)
		builder.WriteString("\treturn []string{}\n")
		builder.WriteString("}\n\n")
	}

	if !methodExists(data.Methods, data.StructName, "FieldsPointers") {
		// Генерация метода FieldsPointers
		fmt.Fprintf(builder, "func (%s *%s) FieldsPointers() []interface{} {\n", receiver, data.StructName)
		builder.WriteString("\treturn []interface{}{\n")
		for _, field := range data.Fields {
			fmt.Fprintf(builder, "\t\t&%s.%s,\n", receiver, field.Name)
		}
		builder.WriteString("\t}\n")
		builder.WriteString("}\n\n")
	}

	return builder.String()
}

// methodExists функция проверки существования метода
func methodExists(methods []MethodInfo, structName, methodName string) bool {
	for _, m := range methods {
		if m.Receiver == structName && m.Name == methodName {
			return true
		}
	}

	return false
}

// AppendToFile функция добавления содержимого в файл
func AppendToFile(filename, content string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		return err
	}
	return nil
}
