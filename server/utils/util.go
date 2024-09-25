package utils

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
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

func SetCookie(ctx *gin.Context, loggedInUserValue *models.User) {
	tokenString, err := CreateJwtToken(loggedInUserValue)
	if err != nil {
		fmt.Println("failed to generate token")
		return
	}

	cookie := http.Cookie{
		Name:     "SessionToken",
		Value:    tokenString,
		Path:     "/",
		MaxAge:   1800,
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(ctx.Writer, &cookie)
}

func GetCookie(ctx *gin.Context) (string, string) {
	cookie, err := ctx.Request.Cookie("SessionToken")
	if err == http.ErrNoCookie {
		return "", ""
	} else if err != nil {
		log.Println("Error: ", err)
		log.Println("Description: Error while Fetching Cookie")
		return "", ""
	}

	tokenString := cookie.Value
	fmt.Println("token string from getcookie func = ", tokenString)

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

func DeleteCookie(ctx *gin.Context) {
	cookie := &http.Cookie{
		Name:   "SessionToken",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(ctx.Writer, cookie)
}

func GetCurrentDate(ctx *gin.Context) string {
	// now := time.Now()
	// fmt.Printf("Current date and time: %s\n", now.Format("2006-01-02 15:04:05"))

	now := time.Now()

	fmt.Println("Current date and time:", now)

	year := now.Year()
	month := now.Month()
	day := now.Day()
	hour := now.Hour()
	minute := now.Minute()
	second := now.Second()

	fmt.Println("Year:", year)
	fmt.Println("Month:", month)
	fmt.Println("Day:", day)
	fmt.Println("Hour:", hour)
	fmt.Println("Minute:", minute)
	fmt.Println("Second:", second)

	date := fmt.Sprintf("%02d-%02d-%04d", now.Day(), now.Month(), now.Year())
	return date
}
