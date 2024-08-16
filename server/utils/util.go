package utils

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/prachin77/server/models"
)

// var resp = make(map[string]string)
// var respInterface = make(map[string]interface{})
var secretKey = []byte("secret-key")

func TokenGenerator() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}


func CreateJwtToken(loggedInUserValue *models.User) (string, error) {
    fmt.Println("Creating token for User:", loggedInUserValue)
    fmt.Println("username in jwt token = ", loggedInUserValue.UserName)
    fmt.Println("user id in jwt token = ", loggedInUserValue.UserId)
    token := jwt.NewWithClaims(jwt.SigningMethodHS256,
        jwt.MapClaims{
            "userid":   loggedInUserValue.UserId,
            "username": loggedInUserValue.UserName,
            "exp":      time.Now().Add(time.Hour * 24).Unix(),
        })

    tokenString, err := token.SignedString(secretKey)
    if err != nil {
        fmt.Println(err)
        return "", err
    }
    fmt.Println("JWT token string : ", tokenString)

    return tokenString, err
}


func SetCookie(w http.ResponseWriter, r *http.Request, loggedInUserValue *models.User) {
	tokenString , err := CreateJwtToken(loggedInUserValue)
	if err!=nil{
		fmt.Println("failed to generate token")
		return
	}

	cookie := http.Cookie{
		Name:     "SessionToken",
		Value:    tokenString,
		Path:     "/",
		MaxAge:   600,
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, &cookie)
}

func GetCookie(w http.ResponseWriter, r *http.Request) (string, string) {
	cookie, err := r.Cookie("SessionToken")

	if err == http.ErrNoCookie {
		return "", ""
	} else if err != nil {
		log.Println("Error: ", err)
		log.Println("Description: Error while Fetching Cookie")
		return "", ""
	}

	tokenString := cookie.Value
	fmt.Println("token string from getcookie func = ",tokenString)

	// Decode JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Make sure the token's signing method matches the expected method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		log.Println("Error parsing token: ", err)
		return "", ""
	}

	// Extract claims from JWT token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userid, ok := claims["userid"].(string)
		if !ok {
			log.Println("Error extracting userid from claims")
			return "", ""
		}
		return userid, tokenString
	} else {
		log.Println("Invalid token")
		return "", ""
	}
}

func DeleteCookie(w http.ResponseWriter, r *http.Request){
	cookie := &http.Cookie{
		Name : "SessionToken",
		Value: "",
		MaxAge: -1,
	}
	http.SetCookie(w,cookie)
}
