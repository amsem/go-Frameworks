package main

import (
	"time"

	"github.com/gin-gonic/gin"
)


func main()  {
    r := gin.Default()
    r.GET("/", func(ctx *gin.Context) {
        ctx.JSON(200, gin.H{
            "servertime": time.Now().UTC(),    
        })
    })
    r.Run(":8000")
}
