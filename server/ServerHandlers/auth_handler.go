package serverhandlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prachin77/db"
	"github.com/prachin77/server/models"
)

func Register(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	var user models.User
	if err := json.NewDecoder(ctx.Request.Body).Decode(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "data not in correct format"})
		return
	}

	if user.UserName == "" || user.Email == "" || user.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "any of the fields should not be empty"})
		return
	}

	userFound, existingUser := db.CheckUserInDB(&user)
	if userFound {
		ctx.JSON(http.StatusConflict, gin.H{"message": "user already exists", "user": existingUser})
		return
	}

	if err := db.CheckUserIdUsingEmail(ctx, &user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "error updating session token (userid)"})
		return
	}

	if err, insertedUser := db.InsertUser(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "error inserting data"})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "registration successful", "user": insertedUser})
	}
}

func Login(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	var user models.User
	if err := json.NewDecoder(ctx.Request.Body).Decode(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "data not in correct format"})
		return
	}

	if user.Email == "" || user.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "email and password can't be empty"})
		return
	}

	userFound, storedUser := db.CheckUserInDB(&user)
	if !userFound {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "no user found in db"})
		return
	}

	if storedUser.UserId == "" {
		// Generate a new session token if UserId is empty
		if err := db.CheckUserIdUsingEmail(ctx, &storedUser); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "error updating session token (userid)"})
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

	if user.UserId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "session token is empty"})
		return
	}

	// Check if the user ID exists in the database
	userExists, err := db.IsTokenPresentInDb(user.UserId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "error checking session token"})
		return
	}

	if !userExists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "invalid session token"})
		return
	}

	if err := db.DeleteUserSessionToken(user.UserId); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "error deleting session token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user logout successful"})
}


