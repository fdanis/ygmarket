package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fdanis/yg-loyalsys/internal/app"
	"github.com/fdanis/yg-loyalsys/internal/db/driver"
	"github.com/fdanis/yg-loyalsys/internal/routes"
)

func main() {
	a := app.NewApp()

	db, err := driver.ConnectSQL(a.Config.ConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	server := &http.Server{
		Addr:    a.Config.Address,
		Handler: routes.Routes(&a, db),
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Printf("server started at %s\n", a.Config.Address)
	<-sig
	log.Print("server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
}
