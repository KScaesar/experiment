# 實驗 jetbrain goland 重構功能

## 情境
legacy code 某個第三方套件的函數  
散落在各個檔案  
現在希望呼叫後, 再多做某些事情

不想要一次手動改多個地方  
所以想說將第三方套件的函數  
提取為自己專案的 package 的函數

原本呼叫 第三方套件函數 的地方  
改為呼叫自己新定義的函數

## 目標
將專案中  
所有如下方形式的程式碼  

```
v:= agollo.GetValue("key")
```

提取變為新的 package 名稱為 config  
新函數為 Get
```
// 原先
v:= agollo.GetValue("key")

---

// 後來
v := config.Get()

func Get() string {
    // do something
    v := agollo.GetValue("key")

    return v
}
```

此專案無法執行

