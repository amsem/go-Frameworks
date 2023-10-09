package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/emicklei/go-restful"
)

func PingTime(req *restful.Request, resp *restful.Response){
    io.WriteString(resp, fmt.Sprintf("%s", time.Now()))
}


func main()  {
    webservice := new(restful.WebService)

    webservice.Route(webservice.GET("/ping").To(PingTime))
    restful.Add(webservice)
    log.Fatal(http.ListenAndServe(":8000", nil))
    
}
