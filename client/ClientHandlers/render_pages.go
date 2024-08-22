package clienthandlers

import (
	"fmt"
	"html/template"

	"github.com/gin-gonic/gin"
	// "github.com/prachin77/server/utils"
)

var (
	PATH = "client/templates/"
)

func DefaultRoute(ctx *gin.Context) {
	ctx.Header("content-Type", "text/html")

	tmpl, err := template.ParseFiles(PATH + "main.html")
	if err != nil {
		fmt.Println("failed to load main.html", err)
		return
	}
	tmpl.Execute(ctx.Writer, nil)
}

func RenderInitPage(ctx *gin.Context) {
	ctx.Header("content-Type", "text/html")
	// userid , tokenString := utils.GetCookie(ctx)
}
