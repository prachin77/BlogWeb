package serverhandlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/prachin77/db"
	"github.com/prachin77/server/models"
	"github.com/prachin77/server/utils"
)

// this response (resp) is to send response message to API Request
var resp = make(map[string]string)
var respInterface = make(map[string]interface{})

func Register(w http.ResponseWriter, r *http.Request) {

	// 1. user registers
	// 2. check whether user with same email & password is already present in db
	// 3. if yes , then prompt user to enter details again
	// 4. if no , then store user details in db
	// 5. set cookies for user in browser

	w.Header().Set("content-Type", "application/json")

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		resp["message"] = "data not in correct format"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if user.UserName == "" || user.Email == "" || user.Password == "" {
		resp["message"] = "any of the fields should not be empty"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	// check if user is present in db or not
	userFound, _ := db.CheckUserInDB(&user)
	if userFound == true {
		resp["message"] = "user found in db !"
		fmt.Println("user found in db !")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Generate Token 
	tokenString := utils.TokenGenerator()
	user.UserId = tokenString		

	err , insertedUser := db.InsertUser(&user)
	if err != nil{
		resp["message"] = "error inserting data"
		fmt.Println("error inserting data")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	fmt.Println("Registration successful")
	fmt.Println("Inserted user details: ", insertedUser)

	respInterface = map[string]interface{}{
		"message": "registration successful",
		"user":    insertedUser,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respInterface)
}

func Login(w http.ResponseWriter, r *http.Request) {

	// 1. user login
	// 2. check user existence through unique password & unique email from DB
	// 3. if login successfull than generate token and assign it to user & update info in DB
	// 4. else give error message & return
	// 5. after token generation create cookie and store its value in browser

	w.Header().Set("content-Type", "application/json")

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil{
		resp["message"] = "data not in correct format"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	// check for empty data
	if user.Email == "" || user.Password == ""{
		resp["message"] = "email and password can't be empty"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	tokenString := utils.TokenGenerator()
	user.UserId = tokenString

	// check through email & pasword if user found  in db or not 
	userFound , storedUser := db.CheckUserInDB(&user)
	if !userFound{
		resp["message"] = "message no user found in db"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return 
	}


	fmt.Println("login sucessfull")
	fmt.Println("login user details = ",storedUser)

	respInterface = map[string]interface{}{
		"message": "registration successful",
		"user":    storedUser,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respInterface)
}

func Logout(w http.ResponseWriter, r *http.Request){
	
}
