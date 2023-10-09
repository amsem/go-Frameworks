package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/amsem/railAPI/dbutils"
	"github.com/emicklei/go-restful"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

type TrainRessource struct{
    ID int 
    DriverName string
    OperatingStatus bool
}

type StationRessource struct {
    ID int
    Name string
    OpeningTime time.Time
    ClosingTime time.Time 
}

type ScheduleRessource struct {
    ID int
    TrainID int
    StationID int
    ArrivalTime time.Time
}

func (t *TrainRessource) Register(container *restful.Container) {
    ws := new(restful.WebService)
    ws.Path("/v1/trains").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
    ws.Route(ws.GET("/{train-id}").To(t.getTrain))
    ws.Route(ws.POST("").To(t.createTrain))
    ws.Route(ws.DELETE("/{train-id}").To(t.removeTrain))
    container.Add(ws)
}

func (t *TrainRessource) getTrain(request *restful.Request, response *restful.Response)  {
    id := request.PathParameter("train-id")
    err := DB.QueryRow("select ID, DRIVER_NAME, OPERATING_STATUS FROM train where id=?", id).Scan(&t.ID, &t.DriverName, &t.OperatingStatus)
    if err != nil {
        log.Println(err)
        response.AddHeader("Content-Type", "text/plain")
        response.WriteErrorString(http.StatusNotFound, "Train could not be found")
    }else {
        response.WriteEntity(t)
    }
}

func (t *TrainRessource) createTrain(request *restful.Request, response *restful.Response)  {
    log.Println(request.Request.Body)
    decoder := json.NewDecoder(request.Request.Body)
    var b TrainRessource
    err := decoder.Decode(&b)
    if err != nil {
        log.Println(err)
    }
    log.Println(b.DriverName, b.OperatingStatus)
    statement, _ := DB.Prepare("insert into train (DRIVER_NAME, OPERATING_STATUS) values (?, ?)")
    result, e := statement.Exec(b.DriverName, b.OperatingStatus)
    if e == nil {
        newID, _ := result.LastInsertId()
        b.ID = int(newID)
        response.WriteHeaderAndEntity(http.StatusCreated, b)
    }else {         
        response.AddHeader("Content-Type", "text/plain")
        response.WriteErrorString(http.StatusInternalServerError, e.Error())
    }
}

func (t *TrainRessource) removeTrain(request *restful.Request, response *restful.Response)  {
    id := request.PathParameter("train-id")
    statement, _ := DB.Prepare("delete from train where id=?")
    _, err := statement.Exec(id)
    if err != nil {
        response.AddHeader("Content-Type", "text/plain")
        response.WriteErrorString(http.StatusInternalServerError, err.Error())
    }else {
        response.WriteHeader(http.StatusOK)
    }
    
}


func main()  {
    var err error
    DB, err = sql.Open("sqlite3", "./railapi.db")
    if err != nil {
        log.Println("Driver creation failed")
    }
    dbutils.Initialize(DB)
    wsContainer := restful.NewContainer()
    wsContainer.Router(restful.CurlyRouter{})
    t := TrainRessource{}
    t.Register(wsContainer)
    log.Println("starting at port 8000")
    server := &http.Server{Addr: ":8000", Handler: wsContainer}
    log.Fatal(server.ListenAndServe())
}
