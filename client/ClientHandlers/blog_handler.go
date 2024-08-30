package clienthandlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prachin77/server/models"
)

func PostBlog(ctx *gin.Context) {
	ctx.Header("content-Type", "text/html")

	var blog models.Blog

	err := json.NewDecoder(ctx.Request.Body).Decode(&blog)
	if err != nil{
		fmt.Println("data not in correct format")
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
}