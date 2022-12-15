package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"net"
	"net/http"
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
		ctx := r.Context()
		select {
		case <-ctx.Done():
			fmt.Println("http graceful shutdown")
			w.WriteHeader(http.StatusOK)
		case <-time.After(2 * time.Second):
			fmt.Fprint(w, "hello world")
		}
	})
	graceful(mux, 8080)
}

func MB(n int) int {
	return n << 20
}

func graceful(h http.Handler, port int) {
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	mainCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
		// Limit the size of the payload.
		Handler: http.MaxBytesHandler(h, int64(MB(1))),
		// Setting timeout is a best practice.
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		// Limit the size of the header.
		MaxHeaderBytes: MB(1),
		BaseContext: func(_ net.Listener) context.Context {
			// https://www.rudderstack.com/blog/implementing-graceful-shutdown-in-go/
			// Pass the mainCtx as the context for every request.
			return mainCtx
		},
	}

	done := make(chan bool)
	go func() {
		// service connections
		log.Printf("(version=%s) listening to port *:8080\n", strings.TrimSpace(version))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
		close(done)
	}()

	<-mainCtx.Done()
	log.Println("shutdown server ...")

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	case <-done:
		log.Println("server successfully terminated")
	}
	log.Println("server exiting")
}
