package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/sing3demons/assessment/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func initDB() *sql.DB {
	connStr := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	createTb := `
	CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);
	`
	_, err = db.Exec(createTb)
	if err != nil {
		log.Fatal("can't create table", err)
	}

	return db
}

func main() {
	db := initDB()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	h := handler.NewApplication(db)

	e.POST("/expenses", h.CreateExpensesHandler)
	// e.GET("/expenses/:id", GetExpensesHandlerByID)
	e.PUT("/expenses/:id", h.UpdateExpensesHandler)
	e.DELETE("/expenses/:id", h.DeleteExpenseHandlerByID)
	// e.GET("/expenses", ListExpensesHandler)
	// e.GET("/expenses", ListExpensesHandler)

	fmt.Println("start at port:", os.Getenv("PORT"))
	go func() {
		if err := e.Start(":" + os.Getenv("PORT")); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown
	fmt.Println("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	fmt.Println("server stop")
}
