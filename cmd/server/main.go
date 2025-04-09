package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"project9/internal/server"
	"time"
)

func main() {

	if len(os.Args) != 3 {
		panic("допустимо только 3 аргумента: <программа> <хост> <dsn>")
	}

	adrs := os.Args[1]
	dsn := os.Args[2]

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(fmt.Errorf("ошибка открытия БД: %w", err))
	}
	defer db.Close()

	pingCtx, cancel := context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	if err = db.PingContext(pingCtx); err != nil {
		panic(fmt.Errorf("ошибка подключения к БД: %w", err))
	}

	server.Start(adrs)
}
