package main

import (
	"bytes"
	"fmt"

	"github.com/spf13/viper"
)

var yamlExample = []byte(`
# --- 表示 文件的開始
# ... 表示 文件的結束
---
Hacker: true
name: steve

# 把下面的 --- and ... 註解, 就可以讀取到後面的值
...
---
age: 35
clothing:
  jacket: leather
...
`)

func main() {
	// 由於不是從檔案讀取, viper 不知道副檔名
	// 所以需要指定解析方式
	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewBuffer(yamlExample))

	// 符號 --- 以下的 key 無法被讀取
	// yaml 的特性, --- 可視為不同檔案
	// 但確實在同一個檔案, 有辦法解決無法讀取的問題嘛?
	key := viper.AllKeys()
	fmt.Println(key)
}
