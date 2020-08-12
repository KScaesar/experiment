package main

import (
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/davecgh/go-spew/spew"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Row struct {
	ID int     `db:"id"`
	C1 *int    `db:"c1"`
	C2 *string `db:"c2"`
	C3 int     `db:"c3"`
	C6 string
	C4 string       `db:"c4"`
	C5 sql.NullTime `db:"c5"`
}

func main() {
	db, err := sqlx.Connect("mysql", "root:1234@/test?parseTime=true&loc=Local")
	if err != nil {
		fmt.Println(err)
	}

	// r := Row{
	// 	C1: nil,
	// 	C2: nil,
	// 	C3: 0,
	// 	C4: "2",
	// 	C5: sql.NullTime{
	// 		Time:  time.Now(),
	// 		Valid: true,
	// 	},
	// }
	// sqlString, args, err := sq.
	// 	Insert("TypeExp").
	// 	Columns("c1", "c2", "c3", "c4", "c5").
	// 	Values(r.C1, r.C2, r.C3, r.C4, r.C5).
	// 	ToSql()

	// result, err := db.Exec(sqlString, args...)
	// if err != nil {
	// 	fmt.Println("db insert", err)
	// }
	// fmt.Println(result.LastInsertId())

	sqlString, args, _ := sq.
		Select("*").
		From("TypeExp").
		Where(sq.Eq{"id": 2}).
		ToSql()
	var data Row
	err = db.Get(&data, sqlString, args...)
	if err != nil {
		fmt.Println("get failed:", err)
	}
	spew.Dump(data)
}

// CREATE TABLE `TypeExp` (
// `id` bigint(20) NOT NULL AUTO_INCREMENT,
// `c1` varchar(100) DEFAULT NULL,
// `c2` int(11) DEFAULT NULL,
// `c3` varchar(100) NOT NULL,
// `c4` int(11) NOT NULL,
// `c5` datetime DEFAULT NULL,
// PRIMARY KEY (`id`)
// ) ENGINE=InnoDB AUTO_INCREMENT=29 DEFAULT CHARSET=utf8mb4
