package serverhandlers

import "github.com/gin-gonic/gin"

func PostBlog(ctx *gin.Context) {
	ctx.Header("content-Type" , "application/json")

}