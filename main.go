package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

type User struct {
	User_id           int       `json:"id"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	Password          string    `json:"password"`
	Registration_date time.Time `json:"registration_date"`
}

var conn *pgxpool.Pool

func init() {
	fmt.Println("Starting the application...")
	godotenv.Load()
	DATABASE_URL := os.Getenv("DATABASE_URL")
	var err error
	conn, err = pgxpool.New(context.Background(), DATABASE_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

}

func main() {
	defer conn.Close()
	e := echo.New()
	e.GET("/user", getUser)
	e.GET("/users", getAllUsers)
	e.POST("/user", handleCreateUser)
	e.POST("/user/update", handleUpdateUser)
	e.Logger.Fatal(e.Start(":1323"))
}

// Database functions
func fetchAllUsers(conn *pgxpool.Pool) []User {
	var users []User
	rows, err := conn.Query(context.Background(), "SELECT * FROM users")
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	for rows.Next() {
		var user User
		err := rows.Scan(&user.User_id, &user.Username, &user.Email, &user.Password, &user.Registration_date)
		if err != nil {
			fmt.Println(err)
		}
		users = append(users, user)
	}

	return users
}

func fetchUser(conn *pgxpool.Pool, user_id string) User {
	row, err := conn.Query(context.Background(), "SELECT * FROM users WHERE user_id = $1", user_id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	row.Next()

	var user User
	row.Scan(&user.User_id, &user.Username, &user.Email, &user.Password, &user.Registration_date)
	if err != nil {
		fmt.Println(err)
	}

	return user
}

func createUser(conn *pgxpool.Pool, user User) bool {
	_, err := conn.Exec(context.Background(), "INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", user.Username, user.Email, user.Password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return false
	}
	return true
}

func updateUser(conn *pgxpool.Pool, user User) bool {
	_, err := conn.Exec(context.Background(), "UPDATE users SET username = $1, email = $2, password = $3 WHERE user_id = $4", user.Username, user.Email, user.Password, user.User_id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return false
	}
	return true
}

// Handler functions
func getUser(c echo.Context) error {
	user_id := c.QueryParam("user_id")
	user := fetchUser(conn, user_id)

	return c.JSON(http.StatusOK, user)
}

func getAllUsers(c echo.Context) error {
	users := fetchAllUsers(conn)

	return c.JSON(http.StatusOK, users)
}

func handleCreateUser(c echo.Context) error {
	user := new(User)
	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&user)
	if err != nil {
		log.Printf("Failed processing in handleCreateUser request: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	wasSucessful := createUser(conn, *user)

	return c.String(http.StatusOK, "The user was created successfully: "+fmt.Sprint(wasSucessful)+"")
}

func handleUpdateUser(c echo.Context) error {
	user := new(User)
	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&user)
	if err != nil {
		log.Printf("Failed processing in handleUpdateUser request: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	wasSucessful := updateUser(conn, *user)

	return c.String(http.StatusOK, "The user was updated successfully: "+fmt.Sprint(wasSucessful)+"")
}
