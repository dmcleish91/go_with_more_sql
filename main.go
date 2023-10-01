package main

import (
	"context"
	"fmt"
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

func init() {
	fmt.Println("Starting the application...")
}

func main() {
	e := echo.New()
	e.GET("/users", getPerson)
	e.Logger.Fatal(e.Start(":1323"))
}

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

func getPerson(c echo.Context) error {
	godotenv.Load()
	DATABASE_URL := os.Getenv("DATABASE_URL")

	conn, err := pgxpool.New(context.Background(), DATABASE_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	users := fetchAllUsers(conn)

	user := User{
		User_id:           users[0].User_id,
		Username:          users[0].Username,
		Email:             users[0].Email,
		Password:          users[0].Password,
		Registration_date: users[0].Registration_date,
	}

	return c.JSON(http.StatusOK, user)
}
