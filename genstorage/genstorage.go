package genstorage

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

// пути до шаблонов
const (
	storageTemplate   = "./genstorage/templates/storageTemplate.tmpl"
	interfaceTemplate = "./genstorage/templates/interfaceTemplate.tmpl"
)

// NewStorage конструктор
func NewStorage(fileName, outputDir string) (*Storage, error) {
	var (
		tableName, structName string
		err                   error
	)

	directory := outputDir
	if outputDir[len(outputDir)-1] != '/' {
		directory += "/"
	}
	if outputDir[0] != '.' || outputDir[1] != '/' {
		directory = "./storage/"
	}

	tableName, err = GetTableName(fileName)
	if err != nil {
		return nil, fmt.Errorf("GetTableName error %v", err)
	}
	if tableName == "" {
		log.Println("Table name not found, use base model")
		tableName = "base"
		structName = "BaseDTO"
		fileName = "base"
	} else {
		structName, err = GetStructName(fileName)
		if err != nil {
			return nil, fmt.Errorf("GetStructName error %v", err)
		}
	}
	// выделяем имя файла
	fileName = strings.TrimSuffix(path.Base(fileName), ".go")

	// убираем возможные подчеркивания, преобразуем первую букву каждого слова tableName в верхний регистр и удаляем пробелы
	formattedTableName := strings.ReplaceAll(strings.Title(strings.ReplaceAll(tableName, "_", " ")), " ", "")

	// переводим tableName в нижний регистр для заполнения шаблона
	tableNameLowercase := strings.ToLower(formattedTableName)

	// выделяем первую букву tableName для заполнения шаблона
	firstLetter := string(tableNameLowercase[0])

	// ищем файл go.mod для дальнейшего парсинга
	rawData, err := SearchFile("go.mod")
	if err != nil {
		return nil, err
	}

	// выделяем модульную строку
	moduleLine, err := ExtractModuleLine(rawData)
	if err != nil {
		return nil, err
	}

	// создаем шаблоны
	storageTemplate, err := NewStorageTemplate()
	if err != nil {
		return nil, err
	}
	interfaceTemplate, err := NewInterfaceTemplate()
	if err != nil {
		return nil, err
	}

	return &Storage{
		FileName:          fileName,
		OutputDir:         directory,
		StorageTemplate:   storageTemplate,
		InterfaceTemplate: interfaceTemplate,
		TemplateData: TemplateData{
			PackageName:         moduleLine,
			TableName:           tableName,
			EntityName:          structName,
			EntityNameLowercase: tableNameLowercase,
			EntityNameUppercase: formattedTableName,
			EntityFirstLetter:   firstLetter,
		},
	}, nil
}

// printHelp функция вывода справки
func printHelp() {
	fmt.Println("Usage:")
	fmt.Println("  app -h           Show help")
	fmt.Println("  app --entity=<file> --output=<directory>")
	fmt.Println()
	fmt.Println("Flags:")
	flag.PrintDefaults()
}

// TemplateData структура с данными для заполнения шаблона
type TemplateData struct {
	PackageName         string // название пакета
	TableName           string // имя таблицы
	EntityName          string // название структуры
	EntityNameLowercase string // название структуры в нижнем регистре
	EntityNameUppercase string // название структуры с большой буквы
	EntityFirstLetter   string // первая буква имени структуры
}

// Storage структура с данными для работы с шаблоном
type Storage struct {
	FileName          string
	OutputDir         string
	StorageTemplate   *template.Template
	InterfaceTemplate *template.Template
	TemplateData      TemplateData
}

// GetTableName функция парсинга имени таблицы из метода TableName в файле
func GetTableName(fileName string) (string, error) {
	// создаем новый набор файлов
	fs := token.NewFileSet()
	// анализируем файл и создаем AST
	node, err := parser.ParseFile(fs, fileName, nil, parser.ParseComments)
	if err != nil {
		return "", err
	}

	var TableName string
	// рекурсивно обходим AST
	ast.Inspect(
		node, func(node ast.Node) bool {
			// определяем тип интерфейсной переменной чтобы затем выполнить код, специфический для этого типа
			switch x := node.(type) {
			// кейс структуры объявления метода
			case *ast.FuncDecl:
				// проверяем, является ли это методом TableName
				if x.Name.Name == "TableName" {
					// проверяем, что первый элемент тела функции это return
					retStmt, ok := x.Body.List[0].(*ast.ReturnStmt)
					if !ok {
						return false
					}
					// проверяем что возвращаемое значение является строкой
					lit, ok := retStmt.Results[0].(*ast.BasicLit)
					if !ok {
						return false
					}
					// убираем кавычки
					TableName = lit.Value[1 : len(lit.Value)-1]
				}
			}
			return true
		},
	)

	return TableName, nil
}

// GetStructName функция парсинга имени структуры
func GetStructName(fileName string) (string, error) {
	// создаем набор файлов для позиционной информации
	fs := token.NewFileSet()
	// создаем AST и анализируем файл
	node, err := parser.ParseFile(fs, fileName, nil, parser.ParseComments)
	if err != nil {
		return "", err
	}

	var structName string
	// рекурсивный обход дерева
	ast.Inspect(
		node, func(node ast.Node) bool {
			// проверяем, является ли текущий узел объявлением типа
			typeSpec, ok := node.(*ast.TypeSpec)
			if !ok {
				return true
			}
			// проверка, является ли текущий тип структурой
			_, ok = typeSpec.Type.(*ast.StructType)
			if ok && structName == "" {
				structName = typeSpec.Name.Name
			}
			return true
		},
	)

	if structName == "" {
		log.Println("could not find struct")
		return "", nil
	}

	return structName, nil
}

// SearchFile функция поиска по файлу
func SearchFile(confName string) ([]byte, error) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println(wd)
	var coursePath string
	coursePath = wd
	courseConfPath := filepath.Join(coursePath, confName)
	for {
		if len(strings.Split(coursePath, "/")) < 2 {
			return nil, fmt.Errorf("%s not found reached: %s", confName, coursePath)
		}
		if _, err = os.Stat(courseConfPath); os.IsNotExist(err) {
			coursePath = filepath.Dir(coursePath)
			courseConfPath = filepath.Join(coursePath, confName)
			continue
		}
		break
	}
	var rawData []byte
	rawData, err = os.ReadFile(courseConfPath)
	if err != nil {
		return nil, fmt.Errorf("cannot open %s path: %s", confName, courseConfPath)
	}

	return rawData, nil
}

// ExtractModuleLine парсит модульную строку из файла go.mod
func ExtractModuleLine(rawData []byte) (string, error) {
	reader := bytes.NewReader(rawData)
	bufReader := bufio.NewReader(reader)

	for {
		line, _, err := bufReader.ReadLine()
		if err == io.EOF {
			break
		}

		if err != nil {
			return "", fmt.Errorf("failed to read go.mod file: %s", err.Error())
		}

		if bytes.Contains(line, []byte("module")) {
			moduleLine := string(line)
			moduleLine = string(bytes.TrimSpace(bytes.TrimPrefix(line, []byte("module"))))
			return moduleLine, nil
		}
	}

	return "", fmt.Errorf("failed to find module line in go.mod file")
}

// NewStorageTemplate конструктор storage шаблона
func NewStorageTemplate() (*template.Template, error) {
	tmpl, err := template.ParseFiles(storageTemplate)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

// NewInterfaceTemplate конструктор interface шаблона
func NewInterfaceTemplate() (*template.Template, error) {
	tmpl, err := template.ParseFiles(interfaceTemplate)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

// CreateStorageFiles функция создания файлов в заданной директории
func (s *Storage) CreateStorageFiles() error {

	// cоздаем директорию
	err := os.MkdirAll(s.OutputDir, 0755)
	if err != nil {
		return err
	}

	files := []struct {
		path     string
		template *template.Template
	}{
		{filepath.Join(s.OutputDir, s.FileName+"_storage.go"), s.StorageTemplate},
		{filepath.Join(s.OutputDir, s.FileName+"_interface.go"), s.InterfaceTemplate},
	}

	for _, f := range files {
		var file *os.File
		// проверяем файл на существование
		if _, err := os.Stat(f.path); os.IsNotExist(err) {
			//если IsNotExist true, создаем файл
			file, err = os.Create(f.path)
			if err != nil {
				return err
			}
			defer file.Close()
			log.Printf("File `%s` created", f.path)
		} else {
			return fmt.Errorf("file `%s` already exists", f.path)
		}

		// заполняем шаблоны
		err := f.template.Execute(file, s.TemplateData)
		if err != nil {
			return err
		}
	}

	return nil
}
