package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/amsem/railAPI-gin/dbutils"
	"github.com/gin-gonic/gin"
	_"github.com/mattn/go-sqlite3"
)



var DB *sql.DB

type StationRessource struct {
    ID int `json:"id"`
    Name string `json:"name"`
    OpeningTime string `json:"opening_time"`
    ClosingTime string `json:"closing_time"`
}

func GetStation(c *gin.Context)  {
    var station StationRessource
    id := c.Param("station_id")
    err := DB.QueryRow("select ID, NAME, CAST(OPENING_TIME as CHAR), CAST(CLOSING_TIME as CHAR) from station where id=?",id).Scan(&station.ID, &station.Name, &station.OpeningTime, &station.ClosingTime)
    if err != nil {
        log.Println("station doesnt exists")
        c.JSON(500, gin.H{
            "error": err.Error(),
        })
    }else {
        c.JSON(200, gin.H{
            "result": station,
        })
    }
}

func CreateStation(c *gin.Context)  {
    var station StationRessource
    if err := c.BindJSON(&station); err == nil {
        statement, _ := DB.Prepare("insert into station (NAME, OPENING_TIME, CLOSING_TIME) values (?, ?, ?)")
        result, _ := statement.Exec(station.Name, station.OpeningTime, station.ClosingTime)
        if err == nil {
        var station StationRessource
            newID, _ := result.LastInsertId()
            station.ID = int(newID)
            c.JSON(http.StatusOK, gin.H{
                "result": station,
            })
        }else {
            c.String(http.StatusInternalServerError, err.Error())
        }
    }else {
        c.String(http.StatusInternalServerError, err.Error())
    }     
}

func RemoveStation(c *gin.Context)  { 
    id := c.Param("station_id")
    statement, _ := DB.Prepare("delete from station where id=?")
    _, err := statement.Exec(id)
    if err != nil {
        log.Println(err)
        c.JSON(500, gin.H{
            "error": err.Error(),
        })
    }else {
        c.String(http.StatusOK, "succes deleting")
    }
}

func main()  {
    var err error
    DB, err = sql.Open("sqlite3", "./rail.db")
    if err != nil {
        log.Println("Driver Creation failed !")
    }
    dbutils.Initialize(DB)
    r := gin.Default()
    r.GET("/v1/stations/:station_id", GetStation)
    r.POST("/v1/stations", CreateStation)
    r.DELETE("/v1/stations/:station_id", RemoveStation)
    r.Run(":8000")
}

