package clienthandlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prachin77/db"
	"github.com/prachin77/server/models"
	"github.com/prachin77/server/utils"
)

// this response (resp) is to send response message to API Request
var resp = make(map[string]string)
var respInterface = make(map[string]interface{})

func Login(ctx *gin.Context) {

	ctx.Header("content-Type", "text/html")

	var user models.User
	err := json.NewDecoder(ctx.Request.Body).Decode(&user)
	if err != nil {
		fmt.Println("data not in correct format")
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// check for empty data
	if user.Email == "" || user.Password == "" {
		fmt.Println("email & password can't be empty")
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// check through email & pasword if user found  in db or not
	userFound, storedUser := db.CheckUserInDB(&user)
	if !userFound {
		fmt.Println("message no user found in db")
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	utils.SetCookie(ctx, &storedUser)

	fmt.Println("login sucessfull")
	fmt.Println("login user details = ", storedUser)

	respInterface = map[string]interface{}{
		"message": "Login successful",
		"user":    storedUser,
	}
	fmt.Println(respInterface)
	ctx.Writer.WriteHeader(http.StatusOK)
	RenderHomePage(ctx,user.UserId)
}

func Register(ctx *gin.Context) {
	ctx.Header("content-Type", "text/html")

	var user models.User
	err := json.NewDecoder(ctx.Request.Body).Decode(&user)
	if err != nil {
		fmt.Println(err)
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if user.UserName == "" || user.Email == "" || user.Password == "" {
		fmt.Println("any of the fields can't be empty")
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// check if user is present in db or not
	userFound, _ := db.CheckUserInDB(&user)
	if userFound == true {
		resp["message"] = "user found in db !"
		ctx.Writer.WriteHeader(http.StatusConflict)
		return
	}

	tokenString := utils.TokenGenerator()
	// Checking whether the generated session token is already used or not
	// If so then generate another session token & keep this loop until we find a unique session token
	for {
		isTokenFound, err := db.IsTokenPresentInDb(tokenString)
		if err != nil {
			fmt.Println(err)
			fmt.Println("can't validate token")
		}
		if !isTokenFound {
			break
		}
		tokenString = utils.TokenGenerator()
		fmt.Println("new session token : ", tokenString)
	}
	user.UserId = tokenString

	err, insertedUser := db.InsertUser(&user)
	if err != nil {
		fmt.Println("error inserting data")
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	utils.SetCookie(ctx, &insertedUser)
	fmt.Println("Registration successful")
	fmt.Println("Inserted user details: ", insertedUser)
	ctx.Writer.WriteHeader(http.StatusOK)
	RenderHomePage(ctx,user.UserId)
}
