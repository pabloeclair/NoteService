package server

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Start(adrs string, dns string) {

	db, err := sql.Open("pgx", dns)
	if err != nil {
		panic(fmt.Errorf("ошибка открытия БД: %w", err))
	}
	defer db.Close()

	pingCtx, cancel := context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	if err = db.PingContext(pingCtx); err != nil {
		panic(fmt.Errorf("ошибка подключения к БД: %w", err))
	}

}
