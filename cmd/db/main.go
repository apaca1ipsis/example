package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"log/slog"

	"database/sql"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

type Handler struct {
	db *sql.DB
}

// github.com/kelseyhightower/envconfig- !!!!!!
func dsn() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("POSTGRES_DB_HOST"),
		os.Getenv("POSTGRES_DB_PORT"),
		os.Getenv("POSTGRES_DB_USER"),
		os.Getenv("POSTGRES_DB_PASSWORD"),
		os.Getenv("POSTGRES_DB_NAME"),
		os.Getenv("POSTGRES_SSL_MODE"))
}

func main() {
	listenPort := os.Getenv("LISTEN_PORT")
	// dsn := os.Getenv("POSTGRES_DSN")
	db, err := sql.Open("postgres", dsn())
	if err != nil {
		panic("db unavailable")
	}
	defer db.Close()

	migrate(db)

	handler := &Handler{db: db}

	fmt.Println("in db/main, port is ", listenPort)
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	fmt.Println("new fiber")

	// example:
	// curl -X POST http://127.0.0.1:9999/api/write/10
	app.Post("/api/write/:number", handler.writeNumber)
	fmt.Println("posted url")

	if err := app.Listen(":" + listenPort); err != nil {
		panic("http server: failed to listen and serve: " + err.Error())
	}

	slog.Info("server started at ", listenPort)
	fmt.Println("server started at ", listenPort)

}

func migrate(db *sql.DB) {
	_, err := db.Exec(
		`create table if not exists logs(
      "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
      "message" varchar not null
    )
     `)
	if err != nil {
		panic(err)
	}
}

var dbLock = sync.Mutex{}

func (h *Handler) writeNumber(c *fiber.Ctx) error {
	dbLock.Lock()
	defer dbLock.Unlock()

	time.Sleep(1 * time.Second)

	number := c.Params("number")

	slog.Info("written",
		"number", number,
	)

	_, err := h.db.Exec(`insert into logs (message) values ($1)`, number)
	if err != nil {
		return err
	}

	return nil
}
