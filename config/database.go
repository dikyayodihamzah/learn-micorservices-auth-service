package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	_ "github.com/joho/godotenv/autoload"
)

var (
	host = os.Getenv("DB_HOST")
	port = os.Getenv("DB_PORT")
	username = os.Getenv("DB_USERNAME")
	password = os.Getenv("DB_PASSWORD")
)

func NewDB(dbName string) *pgx.Conn {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", username, password, host, port, dbName)

	db, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Println(err)
	}

	return db
}