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
	query := "SELECT username, password, session_token FROM users WHERE username = ?"

	row := db.QueryRow(query, user.UserName)

	var storedUser models.User

	err := row.Scan(&storedUser.UserName, &storedUser.Password, &storedUser.SessionToken)
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
	query := "INSERT INTO users (username, password, session_token) VALUES (?, ?, ?)"
	_, err := db.Exec(query, user.UserName, user.Password, user.SessionToken.String)
	if err != nil {
		log.Printf("Error inserting user into database: %v", err)
		return err, models.User{}
	}

	return nil, *user
}


func IsTokenPresentInDb(tokenString string) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE session_token = ?"
	var count int
	err := db.QueryRow(query, tokenString).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func DeleteSessionToken(sessionToken string) error {
    query := "update users set session_token = NULL where session_token = ?"

    _, err := db.Exec(query, sessionToken)
    if err != nil {
        return err
    }

    return nil
}

func UpdateSessionToken(username string, sessionToken string) error {
	query := "UPDATE users SET session_token = ? WHERE username = ?"

	result, err := db.Exec(query, sessionToken, username)
	if err != nil {
		return fmt.Errorf("failed to execute update: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no user found with username: %s", username)
	}

	return nil
}
