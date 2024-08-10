package clienthandlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/prachin77/server/models"
	"github.com/prachin77/server/utils"
)

// "/BlogWeb/client/templates/"++
func DefaultRoute(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("/BlogWeb/client/templates/main.html"))
	tmpl.Execute(w, nil)
}
func RenderHomePage(w http.ResponseWriter, r *http.Request) {
	// cookieValue := utils.GetCookie(w,r)
	// if cookieValue == ""{
	// 	fmt.Println("cookie value is null")
	// 	RenderAuthPage(w,r)
	// }else{
	// 	tmpl := template.Must(template.ParseFiles("/BlogWeb/client/templates/home_page.html"))
	// 	tmpl.Execute(w, nil)
	// }
	userid, tokenString := utils.GetCookie(w, r)

	if tokenString == "" {
		fmt.Println("Cookie value is null")
		RenderLoginPage(w, r)
	} else {
		fmt.Println("UserID from token:", userid)

		// once we recive userid make a get request to home url with userid as param
		// url := fmt.Sprintf("http://localhost:1234/home?userid=%d",userid)
		// fmt.Println("url = ",url)
		// resp , err := http.Get(url)
		// if err != nil{
		// 	fmt.Println(err)
		// }
		// defer resp.Body.Close()

		tmpl := template.Must(template.ParseFiles("/BlogWeb/client/templates/home_page.html"))
		tmpl.Execute(w, userid)
	}
}

func RenderHomeWebPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("/BlogWeb/client/templates/home_page.html"))
	tmpl.Execute(w, nil)
}

func RenderRegisterPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "text/html")
	
    authPageStatus := models.AuthPageStatus{
        IsLogin: false,
    }

    tmpl := template.Must(template.ParseFiles("/BlogWeb/client/templates/auth_page.html"))
    tmpl.Execute(w, authPageStatus)
}


func RenderLoginPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "text/html")

    authPageStatus := models.AuthPageStatus{
        IsLogin: true,
    }

    tmpl := template.Must(template.ParseFiles("/BlogWeb/client/templates/auth_page.html"))
    tmpl.Execute(w, authPageStatus)
}

