package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"net/http"
)

var db *sql.DB
var err error

func main() {
	//docker run --name mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root -d mysql
	db, err = sql.Open("mysql", "root:root@(0.0.0.0:3306)/go?charset=utf8")
	check(err)
	defer db.Close()
	err = db.Ping()
	check(err)

	http.HandleFunc("/", index)
	http.HandleFunc("/amigos", amigos)
	http.HandleFunc("/read", read)
	http.HandleFunc("/create", create)
	http.HandleFunc("/drop", drop)
	http.HandleFunc("/insert", insert)
	http.HandleFunc("/update", update)
	http.HandleFunc("/delete", del)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hi! DB is OK!")
}

func amigos(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`SELECT Name FROM amigos;`)
	check(err)
	defer rows.Close()

	// data to be used in query
	var s, name string
	s = "RETRIEVED RECORDS:\n"

	// query
	for rows.Next() {
		err = rows.Scan(&name)
		check(err)
		s += name + "\n"
	}
	fmt.Fprintln(w, s)
}

func create(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Prepare(`CREATE TABLE customer (Name varchar(255));`)
	check(err)

	res, err := stmt.Exec()
	check(err)

	n, err := res.RowsAffected()
	check(err)

	fmt.Fprintln(w, "Created TABLE customer", n)
}

func drop(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Prepare(`DROP TABLE customer;`)
	check(err)

	res, err := stmt.Exec()
	check(err)

	n, err := res.RowsAffected()

	fmt.Fprintln(w, "Deleted TABLE customer", n)
}

func insert(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Prepare(`INSERT INTO customer VALUES ('Josh');`)
	check(err)

	res, err := stmt.Exec()
	check(err)

	n, err := res.RowsAffected()
	fmt.Fprintln(w, "Inserted in TABLE customer rows", n)
}

func read(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`SELECT * FROM customer;`)
	check(err)

	var s, name string
	s = "SELECT RESULTS:\n"

	for rows.Next() {
		err = rows.Scan(&name)
		check(err)
		s += name + "\n"
	}

	fmt.Fprintln(w, s)
}

func update(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Prepare(`UPDATE customer SET Name = 'Jonny' WHERE Name='Josh';`)
	check(err)

	res, err := stmt.Exec()
	check(err)

	n, err := res.RowsAffected()
	fmt.Fprintln(w, "Updated TABLE customer", n)
}

func del(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Prepare(`DELETE FROM customer WHERE Name='Jonny';`)
	check(err)

	res, err := stmt.Exec()
	check(err)

	n, err := res.RowsAffected()
	fmt.Fprintln(w, "Deleted FROM TABLE customer", n)
}
