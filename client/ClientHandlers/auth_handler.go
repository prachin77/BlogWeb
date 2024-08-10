package clienthandlers

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

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "application/json")

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		resp["message"] = "data might not be in json format !"
		json.NewEncoder(w).Encode(resp)
		return
	}
	fmt.Println("Decoded user from login :", user)

	if user.Email == "" || user.Password == "" {
		resp["message"] = "email or password can't be empty"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	// once we recive value check user presence in db
	userFound, storedUser := db.CheckUserInDB(&user)
	if !userFound {
		resp["message"] = "user not found in db !"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	
	// once we find user in db than generate JWT token 
	// pass logged in user in set cookie 
	utils.SetCookie(w,r,&storedUser) 
	
	fmt.Println("stored user value : ", storedUser)

	respInterface = map[string]interface{}{
		"message": "login successful",
		"user":    storedUser,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respInterface)

}

func Register(w http.ResponseWriter , r* http.Request){
	w.Header().Set("content-type","application/json")

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil{
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		resp["message"] = "data might not be in json format"
		json.NewEncoder(w).Encode(resp)
		return
	}
	fmt.Println("decoder user from register : ",user)

	if user.Email == "" || user.Password == "" || user.UserName == ""{
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		resp["message"] = "username , email or password can't be empty"
		json.NewEncoder(w).Encode(resp)
		return
	}

	// once we recieve user value check user in db 
	userFound , storedUser := db.CheckUserInDB(&user)
	if userFound{
		fmt.Println("user already exists !")
		w.WriteHeader(http.StatusBadRequest)
		resp["message"] = "user already exists !"
		json.NewEncoder(w).Encode(resp)
		return
	}
	fmt.Println("already existing user : ",storedUser)

	// now if we don't find user in db generate new token for user
	tokenString := utils.TokenGenerator()
	fmt.Println("token for new user")
	user.UserId = tokenString

	// now if we don't find user in db store new user in db 
	err , registerdUser := db.InsertUser(&user)
	if err != nil{
		fmt.Println(err)
		resp["message"] = "error registering user !"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	// once we store user db than generate JWT token 
	// pass logged in user in set cookie 
	utils.SetCookie(w,r,&registerdUser) 

	fmt.Println("stored user value : ", storedUser)

	respInterface = map[string]interface{}{
		"message": "registration successful",
		"user":    registerdUser,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respInterface)
}
