package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const url_LM = "https://liquimoly.ru/t_import.php?id=15"

func main(){


	response, err := http.Get(url_LM)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", "./db/lm.db")
	if err != nil {
		log.Fatal(err)
	}

	var responseString = string(responseData)
	responseStrings := strings.Split(responseString, "\n")
	for i := 0; i < len(responseStrings); i++ {
		if responseStrings[i] != "" {
			oneString := strings.Split(responseStrings[i], ";")
			//fmt.Println(oneString)
			insertUrl(db, oneString[0], oneString[1], oneString[2])
		}
	}
}

func insertUrl(db *sql.DB, article string, url string, picture string) {
	insertURLSQL := `INSERT INTO urlinfo(article, url, picture, created) VALUES (?, ?, ?, datetime('now'))`
	statement, err := db.Prepare(insertURLSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(article, url, picture)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
