package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/vaughan0/go-ini"
)

var (
	db  *sql.DB
	err error
)

type dbCredentials struct {
	username string
	password string
	protocol string
	address  string
	dbname   string
}

func ExecuteMultipleRows() {
	rows, err := db.Query("Select actor_id, first_name from sakila.actor limit 5")
	if err != nil {
		fmt.Println("Error Quering", err)
		return
	}
	defer rows.Close()

	var (
		id   int
		name string
	)
	for rows.Next() {
		err = rows.Scan(&id, &name)
		if err != nil {
			fmt.Println("Rows Next", err)
			return
		}
		fmt.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println("Rows.Error", err)
		return
	}
}

func ExecuteSingleRow() {
	var (
		id   int
		name string
	)

	err = db.QueryRow("select actor_id, first_name from actor where actor_id = ?", 1).Scan(&id, &name)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(id, name)
}

func GetDbCredentials() dbCredentials {
	file, err := ini.LoadFile("config.ini")
	if err != nil {
		panic(err)
	}

	lusername, ok := file.Get("mysql", "username")
	if !ok {
		panic("'username' variable missing from 'mysql' section")
	}

	lpassword, ok := file.Get("mysql", "password")
	if !ok {
		panic("'password' variable missing from 'mysql' section")
	}

	lprotocol, ok := file.Get("mysql", "protocol")
	if !ok {
		panic("'protocol' variable missing from 'mysql' section")
	}

	laddress, ok := file.Get("mysql", "address")
	if !ok {
		panic("'address' variable missing from 'mysql' section")
	}

	ldbname, ok := file.Get("mysql", "dbname")
	if !ok {
		panic("'dbname' variable missing from 'mysql' section")
	}

	creds := dbCredentials{
		username: lusername,
		password: lpassword,
		protocol: lprotocol,
		address:  laddress,
		dbname:   ldbname,
	}
	return creds
}

func main() {
	lCreds := GetDbCredentials()
	ldatasource := lCreds.username + ":" + lCreds.password + "@" + lCreds.protocol + "(" + lCreds.address + ")/" + lCreds.dbname
	db, err = sql.Open("mysql", ldatasource)
	if err != nil {
		fmt.Println("Error Opening", err)
		return
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Println("Error Pinging", err)
		return
	}
	ExecuteMultipleRows()
	// ExecuteSingleRow()
}
