package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	clienthandlers "github.com/prachin77/client/ClientHandlers"
)

func main() {
	// Set Gin to release mode
	gin.SetMode(gin.ReleaseMode)

	// Create a new Gin engine without the default middleware
	r := gin.New()

	// Add Logger and Recovery middleware
	r.Use(gin.Logger(), gin.Recovery())

	r.MaxMultipartMemory = 8 << 20 // 8 MiB

	r.GET("/", clienthandlers.DefaultRoute)
	r.GET("/blogger",clienthandlers.RenderInitPage)

	// start server
	err := r.Run(":1234")
	if err != nil {
		fmt.Println(err)
	}
}
