package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Start(adrs string) {

	var shutdownTimeout time.Duration
	var err error

	shutdownTimeoutStr, exists := os.LookupEnv("SHUTDOWN_TIMEOUT")
	if exists {
		shutdownTimeout, err = time.ParseDuration(shutdownTimeoutStr)
		if err != nil {
			fmt.Printf("ошибка валидации SHUTDOWN_TIMEOUT (установелено значение по умолчанию - 7s): %v\n", err)
			shutdownTimeout = 7 * time.Second
		}
	} else {
		shutdownTimeout = 7 * time.Second
	}

	serverCtx, serverCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer serverCancel()

	mux := http.NewServeMux()

	server := http.Server{
		Addr:    adrs,
		Handler: mux,
	}

	log.Printf("Сервер запущен. Адрес: %s. PID: %d\n", adrs, os.Getppid())

	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-serverCtx.Done()
	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancelTimeout()

	server.Shutdown(ctxTimeout)
	log.Println("Сервер закрыт")

}
