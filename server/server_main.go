package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	serverhandlers "github.com/prachin77/server/ServerHandlers"
)

func main() {
	// Set Gin to release mode
	gin.SetMode(gin.ReleaseMode)

	// Create a new Gin engine without the default middleware
	r := gin.New()

	// Add Logger and Recovery middleware
	r.Use(gin.Logger(), gin.Recovery())

	r.MaxMultipartMemory = 8 << 20 // 8 MiB

	// AUTH URLs 
	r.POST("/register", serverhandlers.Register)
	r.POST("/login",serverhandlers.Login)
	r.DELETE("/logout",serverhandlers.Logout)

	// BLOG URLs 
	r.POST("/post-blog",serverhandlers.PostBlog)

	// start server
	err := r.Run(":8080")
	if err != nil {
		fmt.Println(err)
	}
}
