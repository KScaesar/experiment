package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

func main() {
	env, dsn, run := loadEnvFile()
	fmt.Println("env =", env)
	fmt.Println("dsn =", dsn)
	fmt.Println("run =", run)

	pingDatabase(dsn)
	fmt.Println("db ping ok!")
}

func pingDatabase(dsn string) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(fmt.Sprintf("db open: %v", err))
	}

	err = db.Ping()
	if err != nil {
		panic(fmt.Sprintf("db ping: %v", err))
	}

}

func loadEnvFile() (env, dsn, s string) {
	var dir string
	name, exist := os.LookupEnv("ENV")
	if !exist || name == "" {
		dir = "./"
		name = "local"
	} else {
		dir = "/app"
	}

	vp := viper.New()
	vp.SetConfigType("yaml")
	vp.SetConfigName(name)
	vp.AddConfigPath(dir)

	if err := vp.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("load config: %v", err))
	}

	return vp.Get("env").(string), vp.Get("dsn").(string), vp.Get("run").(string)
}
