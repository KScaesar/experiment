package main

import (
	"bytes"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/viper"
)

var yamlExample = []byte(`
repo:
  mysql:
    writer:
      user: "root"
      password: "1234"
      host: "localhost"
      port: "3306"
  redis:
    user: "admin"
    password: "7788"
    address: "localhost:6379"
`)

func main() {
	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewBuffer(yamlExample))

	// viper.Get 可以用類似 fluent style 的方式. 找到對應的 key/value
	fmt.Println(spew.Sdump(viper.Get("repo.mysql.writer")))
	// (map[string]interface {}) (len=4) {
	// 	(string) (len=4) "user": (string) (len=4) "root",
	// 	(string) (len=8) "password": (string) (len=4) "1234",
	// 	(string) (len=4) "host": (string) (len=9) "localhost",
	// 	(string) (len=4) "port": (string) (len=4) "3306"
	// }

	// viper.Unmarshal 無法使用類似 fluent style 的方式來取值
	// (範例如: ProjectConfig.MySQLConfig)
	//
	// 只能一層一層定義相關結構來取值
	// (範例如: ProjectConfig.RedisConfig)
	//
	// 有辦法靠什麼設定值, 達成類似 fluent style ?
	// 還是只能靠自己用反射來實現 ?
	cfg := new(ProjectConfig)
	viper.Unmarshal(cfg)
	fmt.Println(spew.Sdump(cfg))
	// (*main.ProjectConfig)(0xc000101ea0)({
	// 	MySQLWriter: (main.MySQLConfig) {
	// 	 User: (string) "",
	// 	 Password: (string) "",
	// 	 Host: (string) "",
	// 	 Port: (string) ""
	// 	},
	// 	Repo: (struct { Redis main.RedisConfig "mapstructure:\"redis\"" }) {
	// 	 Redis: (main.RedisConfig) {
	// 	  User: (string) (len=5) "admin",
	// 	  Password: (string) (len=4) "7788",
	// 	  Address: (string) (len=14) "localhost:6379"
	// 	 }
	// 	}
	// })
}

type ProjectConfig struct {
	MySQLWriter MySQLConfig `mapstructure:"repo.mysql.writer"`

	Repo struct {
		Redis RedisConfig `mapstructure:"redis"`
	} `mapstructure:"repo"`
}

type RedisConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Address  string `mapstructure:"address"`
}

type MySQLConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
}
