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

	// Serve static files from the "static" directory
	// r.Static("/static", "P:/BlogWeb/client/static/images")
	r.Static("/static", "./static") 

	r.GET("/", clienthandlers.DefaultRoute)
	r.GET("/blogger", clienthandlers.RenderInitPage)

	// AUTH URLs
	r.GET("/register", clienthandlers.RenderRegisterPage)
	r.GET("/login", clienthandlers.RenderLoginPage)
	r.POST("/login", clienthandlers.Login)
	r.POST("/register", clienthandlers.Register)

	// start server
	err := r.Run(":1234")
	if err != nil {
		fmt.Println(err)
	}
}
