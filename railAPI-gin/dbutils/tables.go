package dbutils

import (
	"database/sql"
	"log"
)

func Initialize(dbDriver *sql.DB) {
    statement, driverError := dbDriver.Prepare(train)
    if driverError != nil {
        log.Println(driverError)
    }
    //Create train table
    _, statementError := statement.Exec()
    if statementError != nil {
        log.Println("Table already exists")
    }
    statement, _ = dbDriver.Prepare(station)
    _, statementError = statement.Exec()
    if statementError != nil{
        log.Println("station problem table")
    }
    statement, _ = dbDriver.Prepare(schedule)
    _, statementError = statement.Exec()
    if statementError != nil {
        log.Println("schedule problem")
    }
    log.Println("All tables created successfully")
}
