package genstorage

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	fileName  string
	outputDir string
)

// init вызывается неявно при импорте пакета
func init() {
	flag.StringVar(&fileName, "entity", "", "Path name of entity file")
	flag.StringVar(&outputDir, "output", "./storage/", "Output directory")
	// Создаем новый флаг для вывода справки
	helpFlag := flag.Bool("h", false, "Show help")
	helpLongFlag := flag.Bool("help", false, "Show help")
	flag.Parse()

	// Если установлен флаг "--help" или "-h", выводим справку и завершаем программу
	if *helpFlag || *helpLongFlag {
		printHelp()
		os.Exit(0)
	}
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

// GetFileName функция парсинга имени файла из флага
func GetFileName() (string, error) {
	if fileName == "" {
		flag.PrintDefaults()
		return "", fmt.Errorf("empty flag")
	}
	return fileName, nil
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

// NewStorage конструктор
func NewStorage() (*Storage, error) {
	var (
		tableName, structName string
		err                   error
	)

	fileName, err = GetFileName()
	if err != nil {
		return nil, err
	}
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
			TableName:           tableName,
			EntityName:          structName,
			EntityNameLowercase: tableNameLowercase,
			EntityNameUppercase: formattedTableName,
			EntityFirstLetter:   firstLetter,
		},
	}, nil
}

// NewStorageTemplate конструктор storage шаблона
func NewStorageTemplate() (*template.Template, error) {
	tmpl, err := template.ParseFiles("./genstorage/templates/storageTemplate.tmpl")
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

// NewInterfaceTemplate конструктор interface шаблона
func NewInterfaceTemplate() (*template.Template, error) {
	tmpl, err := template.ParseFiles("./genstorage/templates/interfaceTemplate.tmpl")
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
