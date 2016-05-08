package chat

import (
	"log"
	"database/sql"
	_"github.com/go-sql-driver/mysql"

)
var Db *sql.DB

func connect_db()  {

	var err error

	Db,err = sql.Open("mysql","root:@tcp(localhost:3306)/cloth_passport?charset=utf8");
	if err != nil {
		log.Fatalf("Open database error: %s\n", err)
	}
	defer Db.Close()

	err = Db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

