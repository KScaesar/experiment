package main

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func main() {
	sqlString, args, _ := sq.Select("request_id").From("history").Where(sq.Eq{"request_id": "caesr"}).ToSql()
	fmt.Println(sqlString, args)

	db, err := sqlx.Connect("mysql", "root:1234@/test")
	if err != nil {
		fmt.Println(err)
	}
	var data string
	err = db.Get(&data, sqlString, args...)
	if err != nil {
		fmt.Println("select", err)
	}
	fmt.Println("data=", data)

	fmt.Println()

	condition := map[string]interface{}{"request_id": "casar"}
	sqlString, args, _ = sq.Update("history").Set("note", "love").Where(condition).ToSql()
	fmt.Println(sqlString, args)
	result, err := db.Exec(sqlString, args...)
	if err != nil {
		fmt.Println("update", err)
	}
	affected, _ := result.RowsAffected()
	fmt.Println("RowsAffected", affected)

	updateSql := `
UPDATE
	history
INNER JOIN (
	SELECT
		request_id
	FROM
		history
	where
		request_id = 'caear' ) h 
ON 
	history.request_id = h.request_id 
SET
	history.note = 'atruu'
where
	history.request_id = 'caear'`

	result, err = db.Exec(updateSql)
	if err != nil {
		fmt.Println("updateSql err", err)
	}
	affected, _ = result.RowsAffected()
	fmt.Println("updateSql RowsAffected", affected)
}
