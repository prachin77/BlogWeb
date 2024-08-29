package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/prachin77/server/models"
)

var db *sql.DB

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
	query := "SELECT * FROM users WHERE email = ? AND password = ?"

	row := db.QueryRow(query, user.Email, user.Password)

	var storedUser models.User

	err := row.Scan(&storedUser.UserName, &storedUser.Email, &storedUser.Password, &storedUser.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			// User not found in the database
			return false, models.User{}
		}
		log.Printf("Query execution error: %v", err)
		return false, models.User{}
	}

	// User found in the database
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

func IsTokenPresentInDb(tokenString string) (bool, error) {
    query := "SELECT COUNT(*) FROM users WHERE userid = ?"
    var count int
    err := db.QueryRow(query, tokenString).Scan(&count)
    if err != nil {
        return false, err
    }
    return count > 0, nil
}

func SearchUserWithId(userid string) (models.User,error) {
	query := "SELECT * FROM users WHERE userid = ?"

	row := db.QueryRow(query, userid)

	var storedUser models.User	

	err := row.Scan(&storedUser.UserName , &storedUser.Email , &storedUser.Password , &storedUser.UserId)
	if err != nil{
		if err == sql.ErrNoRows{
			return models.User{} , nil
		}
		return models.User{} ,  nil
	}

	return storedUser , nil
}