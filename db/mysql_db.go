package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/prachin77/server/models"
	"github.com/prachin77/server/utils"
)

var db *sql.DB
var resp = make(map[string]string)

func init() {
	if err := godotenv.Load("P:/BlogWeb/.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	DriverName := os.Getenv("DriverName")
	DataSource := os.Getenv("DataSource")

	var err error
	db, err = sql.Open(DriverName, DataSource)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	fmt.Println("Successfully connected to MySQL")
}

func CheckUserInDB(user *models.User) (bool, models.User) {
	query := "SELECT username, email, password, COALESCE(userid, '') FROM users WHERE email = ? AND password = ?"

	row := db.QueryRow(query, user.Email, user.Password)

	var storedUser models.User

	err := row.Scan(&storedUser.UserName, &storedUser.Email, &storedUser.Password, &storedUser.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, models.User{}
		}
		log.Printf("Query execution error: %v", err)
		return false, models.User{}
	}

	return true, storedUser
}

func InsertUser(user *models.User) (error, models.User) {
	fmt.Println("user details for insert user ")
	fmt.Println(user)
	query := "INSERT INTO users (username, password, email, userid ) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(query, user.UserName, user.Password, user.Email, user.UserId)
	if err != nil {
		log.Printf("Error inserting user into database: %v", err)
		return err, models.User{}
	}

	return nil, *user
}

func SearchUserWithId(userid string) (models.User, error) {
    query := "SELECT username, email, password, userid FROM users WHERE userid = ?"

    row := db.QueryRow(query, userid)

    var user models.User
    err := row.Scan(&user.UserName, &user.Email, &user.Password, &user.UserId)
    if err != nil {
        if err == sql.ErrNoRows {
            // User not found
            return models.User{}, nil
        }
        return models.User{}, err
    }

    return user, nil
}

func IsTokenPresentInDb(tokenString string) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE userid = ?"
	var count int
	err := db.QueryRow(query, tokenString).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func CheckUserIdUsingEmail(ctx *gin.Context, user *models.User) error {
	// If UserId is empty, generate a new one and assign it directly
	if user.UserId == "" {
		user.UserId = utils.TokenGenerator()
		return updateUserIdInDB(ctx, user)
	}

	tokenString := utils.TokenGenerator()

	// Ensure the token is unique
	for {
		isTokenFound, err := IsTokenPresentInDb(tokenString)
		if err != nil {
			return err
		}
		if !isTokenFound {
			break
		}
		tokenString = utils.TokenGenerator()
	}

	user.UserId = tokenString
	return updateUserIdInDB(ctx, user)
}

func updateUserIdInDB(ctx *gin.Context, user *models.User) error {
	insertQuery := "UPDATE users SET userid = ? WHERE email = ?"
	_, err := db.Exec(insertQuery, user.UserId, user.Email)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUserSessionToken(userid string) error {
	query := "UPDATE users SET userid = NULL WHERE userid = ?"

	_, err := db.Exec(query, userid)
	if err != nil {
		log.Printf("Error deleting session token for user ID %s: %v", userid, err)
		return err
	}

	return nil
}
