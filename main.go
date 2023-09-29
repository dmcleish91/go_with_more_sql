package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type User struct {
	user_id           int
	username          string
	email             string
	password          string
	registration_date time.Time
}

func main() {
	godotenv.Load()
	DATABASE_URL := os.Getenv("DATABASE_URL")

	conn, err := pgxpool.New(context.Background(), DATABASE_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fetchAllUsers(conn)
}

func fetchAllUsers(conn *pgxpool.Pool) {
	var users []User
	rows, err := conn.Query(context.Background(), "SELECT * FROM users")
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	for rows.Next() {
		var user User
		err := rows.Scan(&user.user_id, &user.username, &user.email, &user.password, &user.registration_date)
		if err != nil {
			fmt.Println(err)
		}
		users = append(users, user)
	}

	for _, v := range users {
		fmt.Println(v, v.registration_date.GoString())
	}
}
