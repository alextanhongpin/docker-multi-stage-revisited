package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

//go:embed VERSION
var version string

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello world")
	})
	graceful(mux, 8080)
}

func MB(n int) int {
	return n << 20
}

func graceful(h http.Handler, port int) {
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
		// Limit the size of the payload.
		Handler: http.MaxBytesHandler(h, int64(MB(1))),
		// Setting timeout is a best practice.
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		// Limit the size of the header.
		MaxHeaderBytes: MB(1),
	}

	go func() {
		// service connections
		log.Printf("(version=%s) listening to port *:8080\n", strings.TrimSpace(version))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutdown server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("server exiting")
}
