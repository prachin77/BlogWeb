package clienthandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prachin77/db"
	"github.com/prachin77/server/models"
	"github.com/prachin77/server/utils"
)

// this response (resp) is to send response message to API Request
var resp = make(map[string]string)
var respInterface = make(map[string]interface{})

func Register(ctx *gin.Context) {
	ctx.Header("content-Type", "text/html")

	var user models.User
	err := json.NewDecoder(ctx.Request.Body).Decode(&user)
	if err != nil {
		fmt.Println("data not in correct format : ",err)
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if user.UserName == "" || user.Email == "" || user.Password == "" {
		fmt.Println("any of the fields can't be empty")
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	userFound, existingUser := db.CheckUserInDB(&user)
	if userFound {
		resp["message"] = "user found in db !"
		fmt.Println("existing user details : ",existingUser)
		ctx.Writer.WriteHeader(http.StatusConflict)
		return
	}

	if err := db.CheckUserIdUsingEmail(ctx, &user); err != nil {
		fmt.Println("error updating session token (userid)")
		return
	}


	err, insertedUser := db.InsertUser(&user)
	if err != nil {
		fmt.Println("error inserting data")
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.SetCookie(ctx, &insertedUser)
	fmt.Println("Registration successful")
	fmt.Println("Inserted user details: ", insertedUser)
	RenderHomePage(ctx,user.UserId)
	ctx.Writer.WriteHeader(http.StatusOK)
}

func Login(ctx *gin.Context) {
	ctx.Header("content-Type", "text/html")

	var user models.User
	err := json.NewDecoder(ctx.Request.Body).Decode(&user)
	if err != nil {
		fmt.Println("data not in correct format")
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if user.Email == "" || user.Password == "" {
		fmt.Println("email & password can't be empty")
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	userFound, storedUser := db.CheckUserInDB(&user)
	if !userFound {
		fmt.Println("message no user found in db")
		ctx.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	if storedUser.UserId == "" {
		// Generate a new session token if UserId is empty
		if err := db.CheckUserIdUsingEmail(ctx, &storedUser); err != nil {
			fmt.Println("error updating session token (userid)")
			return
		}
	}

	utils.SetCookie(ctx, &storedUser)

	fmt.Println("login sucessfull")
	fmt.Println("login user details = ", storedUser)

	respInterface = map[string]interface{}{
		"message": "Login successful",
		"user":    storedUser,
	}
	fmt.Println(respInterface)
	RenderHomePage(ctx,user.UserId)
	ctx.Writer.WriteHeader(http.StatusOK)
}

func Logout(ctx *gin.Context){
	ctx.Header("content-Type","text/html")

	
	userid, tokenString := utils.GetCookie(ctx)
	if tokenString == "" {
		fmt.Println("cookie is null")
		RenderLoginPage(ctx)
		return
	}
	fmt.Println("cookie token string from logout = ", tokenString)
	fmt.Println("userid from cookie from logout = ", userid)

	// Check if the user ID exists in the database
	userExists, err := db.IsTokenPresentInDb(userid)
	if err != nil {
		fmt.Println("error checking session token")
		return
	}

	if !userExists {
		fmt.Println("invalid session token")
		return
	}

	if err := db.DeleteUserSessionToken(userid); err != nil {
		fmt.Println("error deleting session token")
		return
	}

	utils.DeleteCookie(ctx)
	time.Sleep(5 * time.Second)
	ctx.Writer.WriteHeader(http.StatusOK)
	RenderLoginPage(ctx)
}
