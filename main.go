package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Alexandrhub/cli-orm-gen/genstorage"
)

// флаги командной строки
var (
	fileName  string
	outputDir string
)

// init вызывается неявно при импорте пакета
func init() {
	flag.StringVar(&fileName, "entity", "", "Path name of entity file")
	flag.StringVar(&outputDir, "output", "./repository/", "Output directory")

}

func main() {
	// Создаем новый флаг для вывода справки
	helpFlag := flag.Bool("h", false, "Show help")
	helpLongFlag := flag.Bool("help", false, "Show help")

	flag.Parse()

	if fileName == "" {
		flag.PrintDefaults()
		log.Fatal("empty flag")
	}

	//Если установлен флаг "--help" или "-h", выводим справку и завершаем программу
	if *helpFlag || *helpLongFlag {
		printHelp()
		os.Exit(0)
	}

	// Извлечение информации о структуре
	data, err := genstorage.ReflectFile(fileName)
	if err != nil {
		log.Fatalf("wrong fileName: %v", err)
	}
	// Генерация методов на основе извлеченной информации и вывод результата
	for _, st := range data {
		if !st.HasDBTag {
			continue
		}
		generatedCode := genstorage.GenerateMethods(st)
		if err := genstorage.AppendToFile(fileName, "\n"+generatedCode); err != nil {
			log.Fatalf("AppendToFile error: %v", err)
		}
	}

	// Генерация интерфейса и методов хранилища
	storage, err := genstorage.NewStorage(fileName, outputDir)
	if err != nil {
		log.Fatalf("NewStorage error: %v", err)
	}

	err = storage.CreateStorageFiles()
	if err != nil {
		log.Fatalf("CreateStorageFiles error: %v", err)
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
