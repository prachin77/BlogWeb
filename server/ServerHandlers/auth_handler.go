package serverhandlers

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

func Register(ctx *gin.Context) {

	// 1. user registers
	// 2. check whether user with same email & password is already present in db
	// 3. if yes , then prompt user to enter details again
	// 4. if no , then store user details in db
	// 5. set cookies for user in browser

	ctx.Header("content-Type", "application/json")

	var user models.User
	err := json.NewDecoder(ctx.Request.Body).Decode(&user)
	if err != nil {
		resp["message"] = "data not in correct format"
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(ctx.Writer).Encode(resp)
		return
	}

	if user.UserName == "" || user.Email == "" || user.Password == "" {
		resp["message"] = "any of the fields should not be empty"
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(ctx.Writer).Encode(resp)
		return
	}

	// check if user is present in db or not
	userFound, _ := db.CheckUserInDB(&user)
	if userFound == true {
		resp["message"] = "user found in db !"
		fmt.Println("user found in db !")
		ctx.Writer.WriteHeader(http.StatusConflict)
		json.NewEncoder(ctx.Writer).Encode(resp)
		return
	}

	tokenString := utils.TokenGenerator()

	// Checking whether the generated session token is already used or not
	// If so then generate another session token & keep this loop until we find a unique session token
	for {
		isTokenFound , err := db.IsTokenPresentInDb(tokenString)
		if err != nil{
			fmt.Println(err)
			resp["message"] = "cant validate session token"
			json.NewEncoder(ctx.Writer).Encode(resp)
		}
		if !isTokenFound{
			break
		}
		tokenString = utils.TokenGenerator()
		fmt.Println("new session token : ",tokenString)
	}

	user.UserId = tokenString

	err, insertedUser := db.InsertUser(&user)
	if err != nil {
		resp["message"] = "error inserting data"
		fmt.Println("error inserting data")
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(ctx.Writer).Encode(resp)
		return
	}

	fmt.Println("Registration successful")
	fmt.Println("Inserted user details: ", insertedUser)

	respInterface = map[string]interface{}{
		"message": "registration successful",
		"user":    insertedUser,
	}
	ctx.Writer.WriteHeader(http.StatusOK)
	json.NewEncoder(ctx.Writer).Encode(respInterface)
}

func Login(ctx *gin.Context) {

	// 1. user login
	// 2. check user existence through unique password & unique email from DB
	// 3. if login successfull than generate token and assign it to user & update info in DB
	// 4. else give error message & return
	// 5. after token generation create cookie and store its value in browser

	ctx.Header("content-Type", "application/json")

	var user models.User
	err := json.NewDecoder(ctx.Request.Body).Decode(&user)
	if err != nil {
		resp["message"] = "data not in correct format"
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(ctx.Writer).Encode(resp)
		return
	}

	// check for empty data
	if user.Email == "" || user.Password == "" {
		resp["message"] = "email and password can't be empty"
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(ctx.Writer).Encode(resp)
		return
	}


	// check through email & pasword if user found  in db or not
	userFound, storedUser := db.CheckUserInDB(&user)
	if !userFound {
		resp["message"] = "message no user found in db"
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(ctx.Writer).Encode(resp)
		return
	}

	fmt.Println("login sucessfull")
	fmt.Println("login user details = ", storedUser)

	respInterface = map[string]interface{}{
		"message": "Login successful",
		"user":    storedUser,
	}

	ctx.Writer.WriteHeader(http.StatusOK)
	json.NewEncoder(ctx.Writer).Encode(respInterface)
}
