package main

import (
	"database/sql"
	"log"
    _ "github.com/mattn/go-sqlite3"
)


           
type Book struct {
    id int
    name string
    author string
}

func dbOperations(db *sql.DB)  {
    statement, _ := db.Prepare("INSERT INTO books (name, author, isbn) VALUES (?, ?, ?)")
    statement.Exec("A Tale of Two Cities", "Charles Dickness", 3241512354)
    log.Println("Inserted the book to the db")

    rows, _ := db.Query("SELECT id, name, author FROM books")

    var tempsBook Book
    for rows.Next() {
        rows.Scan(&tempsBook.id, &tempsBook.name, &tempsBook.author)
        log.Printf("ID:%d , Name: %s, Author: %s\n", tempsBook.id, tempsBook.name, tempsBook.author)
    }

    //Update 
    statement, _ = db.Prepare("update books set name=? where id=?")
    statement.Exec("The Tale of Two Cities", 1)
    log.Println("successfully updated the DB")
    //delete
    statement, _ = db.Prepare("delete from books where id=?")

    statement.Exec(1)
    log.Println("successfully deleted the book in db")
}


func main()  {
    db, err := sql.Open("sqlite3", "./books.db")
    if err != nil {
        log.Println(err)
    }
    statement, e := db.Prepare("CREATE TABLE IF NOT EXISTS books(id INTEGER PRIMARY KEY, isbn INTEGER, author VARCHAR(64),name VARCHAR(64) NULL)")
    if e != nil {
        log.Println("Error creating table books")
    }else {
        log.Println("successfully created table books")
    }
    statement.Exec()
    dbOperations(db)
    
}
