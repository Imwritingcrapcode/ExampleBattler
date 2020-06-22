package main

import (
	"database/sql"
	"log"
	"strconv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	DATABASE, err := sql.Open("sqlite3", "Server\\Battler.db")
	if err != nil {
		panic(err)
	}
	st, err := DATABASE.Prepare("UPDATE notifications SET seen = 1 WHERE userID = " + strconv.FormatInt(4, 10) +
		" AND redirect = '" + "friends" + "'")
	if err != nil {
		log.Println("[See notifications]", err.Error())
	} else {
		res, _ := st.Exec()
		b, _ := res.LastInsertId()
		s, _ := res.RowsAffected()
		log.Println(s, b)
	}
}
