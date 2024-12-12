package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strings"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "tan"
	dbname   = "blogdb"
)

var db *sql.DB

func connectDB() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db, nil
}

func getAllCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	getAllCategoriesSQLStatement := "SELECT * FROM categories"
	rows, err := db.Query(getAllCategoriesSQLStatement)
	if err != nil {
		panic(err)
	}
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		panic(err)
	}

	count := len(columnTypes)
	finalRows := []interface{}{}

	for rows.Next() {

		scanArgs := make([]interface{}, count)

		for i, v := range columnTypes {

			switch v.DatabaseTypeName() {
			case "VARCHAR", "TEXT", "UUID", "TIMESTAMP":
				scanArgs[i] = new(sql.NullString)
				break
			case "BOOL":
				scanArgs[i] = new(sql.NullBool)
				break
			case "INT4":
				scanArgs[i] = new(sql.NullInt64)
				break
			default:
				scanArgs[i] = new(sql.NullString)
			}
		}

		err := rows.Scan(scanArgs...)

		if err != nil {
			panic(err)
		}

		masterData := map[string]interface{}{}

		for i, v := range columnTypes {

			if z, ok := (scanArgs[i]).(*sql.NullBool); ok {
				masterData[v.Name()] = z.Bool
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullString); ok {
				masterData[v.Name()] = z.String
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullInt64); ok {
				masterData[v.Name()] = z.Int64
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullFloat64); ok {
				masterData[v.Name()] = z.Float64
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullInt32); ok {
				masterData[v.Name()] = z.Int32
				continue
			}

			masterData[v.Name()] = scanArgs[i]
		}

		finalRows = append(finalRows, masterData)
	}

	categoriesData, err := json.Marshal(finalRows)
	w.Write([]byte(categoriesData))
}

func getAllPostsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	getAllCategoriesSQLStatement := "SELECT * FROM posts"
	rows, err := db.Query(getAllCategoriesSQLStatement)
	if err != nil {
		panic(err)
	}
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		panic(err)
	}

	count := len(columnTypes)
	finalRows := []interface{}{}

	for rows.Next() {

		scanArgs := make([]interface{}, count)

		for i, v := range columnTypes {

			switch v.DatabaseTypeName() {
			case "VARCHAR", "TEXT", "UUID", "TIMESTAMP":
				scanArgs[i] = new(sql.NullString)
				break
			case "BOOL":
				scanArgs[i] = new(sql.NullBool)
				break
			case "INT4":
				scanArgs[i] = new(sql.NullInt64)
				break
			default:
				scanArgs[i] = new(sql.NullString)
			}
		}

		err := rows.Scan(scanArgs...)

		if err != nil {
			panic(err)
		}

		masterData := map[string]interface{}{}

		for i, v := range columnTypes {

			if z, ok := (scanArgs[i]).(*sql.NullBool); ok {
				masterData[v.Name()] = z.Bool
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullString); ok {
				masterData[v.Name()] = z.String
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullInt64); ok {
				masterData[v.Name()] = z.Int64
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullFloat64); ok {
				masterData[v.Name()] = z.Float64
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullInt32); ok {
				masterData[v.Name()] = z.Int32
				continue
			}

			masterData[v.Name()] = scanArgs[i]
		}

		finalRows = append(finalRows, masterData)
	}

	categoriesData, err := json.Marshal(finalRows)
	w.Write([]byte(categoriesData))
}

func getCategoryByIdHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(strings.Split(r.URL.Path, "/"))
}

func getPostByIdHandler(w http.ResponseWriter, r *http.Request) {
}

func addCategoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	categorySqlStatement := `
      INSERT INTO categories(category_name) 
      VALUES ($1)
      RETURNING id
    `
	id := 0
	err := db.QueryRow(categorySqlStatement, "Programming").Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Println("New Category Record ID is: ", id)
	w.Write([]byte(fmt.Sprintf("category %d", id)))
}

func addPostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	postsSqlStatement := `
      INSERT INTO posts (message, imageUrl)
      VALUES ($1, $2)
      RETURNING id
    `
	id := 0
	err := db.QueryRow(postsSqlStatement, "hello", "www.google.com").Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Println("New Post Record ID is: ", id)
	w.Write([]byte(fmt.Sprintf("post %d", id)))
}

func deleteCategoryByIdHandler(w http.ResponseWriter, r *http.Request) {
}

func deletePostByIdHandler(w http.ResponseWriter, r *http.Request) {
}

func main() {
	db, _ = connectDB()
	http.HandleFunc("/categories", getAllCategoriesHandler)
	http.HandleFunc("/posts", getAllPostsHandler)
	http.HandleFunc("/category/{id}", getCategoryByIdHandler)
	http.HandleFunc("/post/:id", getPostByIdHandler)
	http.HandleFunc("/addCategory", addCategoryHandler)
	http.HandleFunc("/addPost", addPostHandler)
	http.HandleFunc("/deleteCategory/:id", deleteCategoryByIdHandler)
	http.HandleFunc("/deletePost/:id", deletePostByIdHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
