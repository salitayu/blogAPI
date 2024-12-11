package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "tan"
	dbname   = "blogdb"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
	categorySqlStatement := `
      INSERT INTO categories(category_name) 
      VALUES ($1)
      RETURNING id
    `
	id := 0
	err = db.QueryRow(categorySqlStatement, "Programming").Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Println("New Category Record ID is: ", id)
	postsSqlStatement := `
      INSERT INTO posts (message, imageUrl)
      VALUES ($1, $2)
      RETURNING id
    `
	err = db.QueryRow(postsSqlStatement, "hello", "www.google.com").Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Println("New Post Record ID is: ", id)
}
