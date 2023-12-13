package adapter

import (
	"reflect"
	"strings"
)

//go:generate go get -u github.com/valyala/quicktemplate/qtc
//go:generate qtc -dir=../../infrastructure/db/migrate

func GetFieldsPointers(u interface{}, args ...string) []interface{} {
	val := reflect.ValueOf(u).Elem()
	v := make([]interface{}, 0, val.NumField())
	
	for i := 0; i < val.NumField(); i++ {
		if len(args) != 0 {
			tagsRaw := val.Type().Field(i).Tag.Get("db_ops")
			tags := strings.Split(tagsRaw, ",")
			found := false
			for _, tag := range tags {
				if tag == args[0] {
					found = true
				}
			}
			if !found {
				continue
			}
		}
		valueField := val.Field(i)
		v = append(v, valueField.Addr().Interface())
	}
	
	return v
}

type In struct {
	Field string
	Args  []interface{}
}
