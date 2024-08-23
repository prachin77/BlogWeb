package clienthandlers

import (
	"fmt"
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/prachin77/server/models"
	"github.com/prachin77/server/utils"
	// "github.com/prachin77/server/utils"
)

var (
	PATH = "P:/BlogWeb/client/templates/"
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
	userid, tokenString := utils.GetCookie(ctx)
	if tokenString == "" {
		fmt.Println("cookie is null")
		RenderLoginPage(ctx)
		return
	}
	fmt.Println("cookie token string = ", tokenString)
	fmt.Println("userid from cookie = ", userid)
	RenderHomePage(ctx,userid)
}

func RenderHomePage(ctx *gin.Context , userid string) {
	ctx.Header("content-Type", "text/html")

	tmpl, err := template.ParseFiles(PATH + "homepage.html")
	if err != nil {
		fmt.Println("unable to render home page = ", err)
		return
	}
	tmpl.Execute(ctx.Writer, userid)
}

func RenderLoginPage(ctx *gin.Context) {
	ctx.Header("content-Type", "text/html")

	tmpl, err := template.ParseFiles(PATH + "auth_page.html")
	if err != nil {
		fmt.Println("unable to render auth page = ", err)
		return
	}
	authStatus := models.AuthPageStatus{
		IsLogin: true,
	}
	tmpl.Execute(ctx.Writer, authStatus)
}

func RenderRegisterPage(ctx *gin.Context) {
	ctx.Header("content-Type", "text/html")
	authStatus := models.AuthPageStatus{
		IsLogin: false,
	}
	tmpl := template.Must(template.ParseFiles(PATH + "auth_page.html"))
	tmpl.Execute(ctx.Writer, authStatus)
}
