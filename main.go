package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/story-pointify/room-service/ws"
)

func main() {
	wsServer := ws.NewWebSocketServer()
	go wsServer.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(wsServer, w, r)
	})

	srv := &http.Server{
		Addr:    ":1337",
		Handler: nil,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Println("Server started 🚀")

	<-done
	log.Println("Server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("💩 Server Shutdown Failed:%v+", err)
	}

	log.Println("Server exited properly 😴")
}
