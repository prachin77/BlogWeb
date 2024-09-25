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

func Register(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	var user models.User
	if err := json.NewDecoder(ctx.Request.Body).Decode(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "data not in correct format"})
		return
	}

	if user.UserName == "" || user.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "any of the fields should not be empty"})
		return
	}

	userFound, existingUser := db.CheckUserInDB(&user)
	if userFound {
		ctx.JSON(http.StatusConflict, gin.H{"message": "user already exists", "user": existingUser})
		return
	}

	for {
		user.SessionToken.String = utils.TokenGenerator()

		isUserIdPresent, err := db.IsTokenPresentInDb(user.SessionToken.String)
		if err != nil {
			fmt.Println("Error checking user ID in DB: ", err)
			return
		}
	
		if !isUserIdPresent {
			break
		}
	}


	if err, insertedUser := db.InsertUser(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "error inserting data"})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "registration successful", "username":insertedUser.UserName,"password":insertedUser.Password,"session_token":insertedUser.SessionToken})
	}
}

func Login(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	var user models.User
	if err := json.NewDecoder(ctx.Request.Body).Decode(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "data not in correct format"})
		return
	}

	if user.UserName == "" || user.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "username and password can't be empty"})
		return
	}

	userFound, storedUser := db.CheckUserInDB(&user)
	if !userFound {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "no user found in db"})
		return
	}

	if storedUser.SessionToken.String == "" {
		for {
			storedUser.SessionToken.String = utils.TokenGenerator()

			isTokenPresent, err := db.IsTokenPresentInDb(storedUser.SessionToken.String)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"message": "error checking session token"})
				return
			}

			if !isTokenPresent {
				break
			}
		}

		if err := db.UpdateSessionToken(storedUser.UserName, storedUser.SessionToken.String); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "error updating session token"})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": storedUser})
}


func Logout(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	var user models.User
	if err := json.NewDecoder(ctx.Request.Body).Decode(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "data not in correct format"})
		return
	}

	if user.SessionToken.String == ""{
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "session token is empty"})
		return
	}

	userExists, err := db.IsTokenPresentInDb(user.SessionToken.String)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "error checking session token"})
		return
	}

	if !userExists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "invalid session token"})
		return
	}

	// Delete session token 
	if err := db.DeleteSessionToken(user.SessionToken.String); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "error deleting session token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user logout successful"})
}

