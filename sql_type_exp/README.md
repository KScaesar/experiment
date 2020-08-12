
## 若 go 為 字串型別, mysql 為 整數型別  
- go 傳空字串, 無法寫入到 mysql  
Error 1366: Incorrect integer value: '' for column 'c4' at row 1

- go 傳 字串 '2', 可順利寫入資料庫, 在 mysql 變為整數 2

- mysql 為數字 2, go 執行 select 可以轉換為 字串 '2'

## 若 go 為 整數型別, mysql 為 字串型別 
- go 傳 整數 0, 可順利寫入資料庫, 在 mysql 變為字串 '0'

- mysql 為字串 '1', go 執行 select 可以轉換為 整數1

- mysql 為空字串 '', go 執行 select 會發生錯誤, 無法轉換為整數  
sql: Scan error on column index 3, name "c3": converting driver.Value type []uint8 ("''") to a int: invalid syntax

## auto_increment  
auto_increment 的欄位, 不可為字串型別  

## 資料庫欄位 可 NULL 值

若 go 用 time.Time 型別來讀取 NULL 值  
會發生錯誤  
sql: Scan error on column index 0, name "c5": unsupported Scan, storing driver.Value type <nil> into type *time.Time

若 go 用 ＊time.Time 型別來讀取 NULL 值  
可正常讀取 

## sqlx.db 接資料

若 select 資料庫所有欄位  
如果 ROW 型別缺少, 返回值的部份欄位  
在使用 db.Get 時, 會發生錯誤  
因為不知道要把回傳值, 塞到哪個欄位  
missing destination name c4 in *main.Row

ROW 型別 想要接資料  
必須包含 select 返回值的所有欄位  
ROW 型別 可以有多餘的欄位, 但不能缺少  
```
sqlString, args, _ := sq.
	Select("*").
	From("TypeExp").
	Where(sq.Eq{"id": 2}).
	ToSql()

var data Row
err = db.Get(&data, sqlString, args...)
```