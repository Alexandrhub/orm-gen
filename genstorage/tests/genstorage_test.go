package tests

import (
	"flag"
	"os"
	"testing"
	"text/template"

	"github.com/Alexandrhub/cli-orm-gen/genstorage"
)

func TestStorage_CreateStorageFiles(t *testing.T) {
	emptyTemplate := template.New("mock")
	mockTemplate := template.New("mock")
	_, err := mockTemplate.Parse("Mock template")
	if err != nil {
		t.Errorf("Error parsing template: %v", err)
	}

	type fields struct {
		Entity            string
		OutputDir         string
		StorageTemplate   *template.Template
		InterfaceTemplate *template.Template
		TemplateData      genstorage.TemplateData
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		setup   func()
	}{
		{
			name: "valid data",
			fields: fields{
				Entity:            "test_model",
				OutputDir:         "./storage/",
				StorageTemplate:   mockTemplate,
				InterfaceTemplate: mockTemplate,
				TemplateData:      genstorage.TemplateData{},
			},
			wantErr: false,
		},
		{
			name: "wrong template",
			fields: fields{
				Entity:            "test_model",
				OutputDir:         "./storage/",
				StorageTemplate:   emptyTemplate,
				InterfaceTemplate: emptyTemplate,
				TemplateData:      genstorage.TemplateData{},
			},
			wantErr: true,
		},
		{
			name: "files already exists",
			fields: fields{
				Entity:            "test_model",
				OutputDir:         "./storage/",
				StorageTemplate:   mockTemplate,
				InterfaceTemplate: mockTemplate,
				TemplateData:      genstorage.TemplateData{},
			},
			setup: func() {
				os.MkdirAll("./storage/", 0o755)
				file, _ := os.Create("./storage/test_model_interface.go")

				file.Close()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				// выполнение настроек перед тестом, если они определены
				if tt.setup != nil {
					tt.setup()

					// удаляем директорию после завершения теста
					defer os.RemoveAll("./storage/")
				}

				// cоздание экземпляра Storage для тестирования
				s := &genstorage.Storage{
					FileName:          tt.fields.Entity,
					OutputDir:         tt.fields.OutputDir,
					StorageTemplate:   tt.fields.StorageTemplate,
					InterfaceTemplate: tt.fields.InterfaceTemplate,
					TemplateData:      tt.fields.TemplateData,
				}

				if err := s.CreateStorageFiles(); (err != nil) != tt.wantErr {
					t.Errorf("CreateStorageFiles() error = %v, wantErr %v", err, tt.wantErr)
				}

				// проверка создания файлов
				if !tt.wantErr {
					for _, filepath := range []string{
						"./storage/" + tt.fields.Entity + "_storage.go",
						"./storage/" + tt.fields.Entity + "_interface.go",
					} {
						if _, err := os.Stat(filepath); os.IsNotExist(err) {
							t.Errorf("File %s is not created", filepath)
						}
					}
				}
			},
		)
	}
}

func TestGetTableName(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "valid method name",
			args:    args{fileName: "test_model.go"},
			want:    "test",
			wantErr: false,
		},
		{
			name:    "not valid method name",
			args:    args{fileName: "wrong_model.go"},
			want:    "",
			wantErr: false,
		},
		{
			name:    "empty model",
			args:    args{fileName: "empty_model.go"},
			want:    "",
			wantErr: false,
		},
		{
			name:    "not valid file name",
			args:    args{fileName: ""},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := genstorage.GetTableName(tt.args.fileName)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetTableName() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("GetTableName() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestGetStructName(t *testing.T) {
	type args struct {
		entity string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "valid struct model",
			args:    args{entity: "test_model.go"},
			want:    "Test",
			wantErr: false,
		},
		{
			name:    "empty struct model",
			args:    args{entity: "empty_model.go"},
			want:    "",
			wantErr: false,
		},
		{
			name:    "not valid file name",
			args:    args{entity: "..."},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := genstorage.GetStructName(tt.args.entity)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetStructName() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("GetStructName() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestGetFileName(t *testing.T) {
	tests := []struct {
		name      string
		flagValue string
		want      string
		wantErr   bool
	}{
		{
			name:      "valid flag",
			flagValue: "testfile.go",
			want:      "testfile.go",
			wantErr:   false,
		},
		{
			name:      "empty flag",
			flagValue: "",
			want:      "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				flag.Set("entity", tt.flagValue)
				got, err := genstorage.GetFileName()
				if (err != nil) != tt.wantErr {
					t.Errorf("GetFileName() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("GetFileName() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
