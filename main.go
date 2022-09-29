package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const url_LM = "https://liquimoly.ru/t_import.php?id=15"

func main() {

	response, err := http.Get(url_LM)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	os.Remove("./db/urlsBrands.db")

	_, err = os.Create("./db/urlsBrands.db") // Create SQLite file
	if err != nil {
		log.Fatal(err.Error())
	}

	db, _ := sql.Open("sqlite3", "./db/urlsBrands.db")
	defer db.Close()
	createTable(db)

	var responseString = string(responseData)
	responseStrings := strings.Split(responseString, "\n")
	for i := 0; i < len(responseStrings); i++ {
		if responseStrings[i] != "" {
			oneString := strings.Split(responseStrings[i], ";")
			//article;url;picture_url
			insertUrl(db, oneString[0], "LM", oneString[1], oneString[2])
		}
	}

	displayURLs(db)
}

func insertUrl(db *sql.DB, article string, brand string, url string, picture string) {
	insertURLSQL := `INSERT INTO urlinfo(article, brand, URL, picture) VALUES (?, ?, ?, ?)`
	statement, err := db.Prepare(insertURLSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(article, brand, url, picture)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func createTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE urlinfo (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"article" TEXT,
		"brand" TEXT,		
		"URL" TEXT,
		"picture" TEXT		
	  );`

	statement, err := db.Prepare(createTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
}

func displayURLs(db *sql.DB) {
	row, err := db.Query("SELECT article, URL, picture FROM urlinfo WHERE brand = ? ORDER BY article", "LM")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() {
		var article string
		var url string
		var picture string
		row.Scan(&article, &url, &picture)
		log.Printf("Article: %s URL: %s Picture: %s", article, url, picture)
	}
}
