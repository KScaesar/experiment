package main

import (
	"fmt"
	"log"
	"reflect"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/viper"
)

func init() {

}

func main() {
	cfg := NewConfigV2("../config", "dev")
	c := spew.Sdump(cfg)
	// fmt.Println(cfg.MySQLReader, cfg.MySQL.Reader)
	fmt.Println(c)
}

func NewConfigV2(path string, fileName string) *ProjectConfig {
	vp := newViper(path, fileName)
	println(spew.Sdump(vp.Get("writer")))

	cfg := new(ProjectConfig)
	structValue := reflect.ValueOf(cfg).Elem()
	parentKey := ""
	SearchAndSet(structValue, vp, parentKey)
	return cfg
}

func SearchAndSet(structValue reflect.Value, vp *viper.Viper, parentKey string) {
	keyTag := "config"
	for i := 0; i < structValue.NumField(); i++ {
		fieldInfo := structValue.Type().Field(i)
		key, ok := fieldInfo.Tag.Lookup(keyTag)
		if !ok {
			continue
		}
		fieldValue := structValue.Field(i)
		if !fieldValue.CanSet() {

		}
		key = parentKey + key
		switch fieldInfo.Type.Kind() {
		case reflect.String:
			v, ok := vp.Get(key).(string)
			if !ok {
				panic("not found key")
			}
			fieldValue.SetString(v)
		case reflect.Bool:
			v := vp.Get(key).(bool)
			fieldValue.SetBool(v)
		case reflect.Struct:
			SearchAndSet(fieldValue, vp, key)
		default:
			panic("not allow type")
		}

	}
}

func newViper(path string, fileName string) *viper.Viper {
	if path == "" || fileName == "" {
	}

	v := viper.New()
	v.SetConfigType("yaml")
	v.AddConfigPath(path)
	v.SetConfigName(fileName)

	// If a config file is found, read it in.
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
		return nil
	}
	log.Println("Using config file:", v.ConfigFileUsed())
	return v
}
