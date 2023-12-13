package main

import (
	"log"
	
	"github.com/Alexandrhub/cli-orm-gen/genstorage"
)

// go build cmd/cli-orm/repogen.go
// ./repogen -entity="./genstorage/models/test_model.go" - создаем директорию с файлами хранилища в корневом каталоге
func main() {
	storage, err := genstorage.NewStorage()
	if err != nil {
		log.Fatal(err)
	}
	
	err = storage.CreateStorageFiles()
	if err != nil {
		log.Fatal(err)
	}
}
