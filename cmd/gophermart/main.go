package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fdanis/yg-loyalsys/internal/accrualclient"
	"github.com/fdanis/yg-loyalsys/internal/app"
	"github.com/fdanis/yg-loyalsys/internal/routes"
)

func main() {
	a := app.NewApp()
	server := &http.Server{
		Addr:    a.Config.Address,
		Handler: routes.Routes(&a),
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	client, err := accrualclient.NewClient(a.Config.AccrualSystemAddress, a.OrderRepository)
	if err != nil {
		log.Fatalf("create accrual client %v \n", err)
	}
	ctxClient, cancelClient := context.WithCancel(context.Background())
	defer cancelClient()
	go runClient(ctxClient, client)

	log.Printf("server started at %s\n", a.Config.Address)
	<-sig
	log.Print("server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
}

func runClient(ctx context.Context, client *accrualclient.Client) {
	t := time.NewTicker(time.Second)
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()

	for {
		select {
		case <-t.C:
			{
				client.Run(ctx2)
			}
		case <-ctx.Done():
			t.Stop()
		}
	}
}
