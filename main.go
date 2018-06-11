package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
)

var (
	logger = log.New(os.Stdout, "allocator: ", log.LstdFlags)
)

var (
	srv  http.Server
	bulk uint64
	used uint64
)

func main() {
	flag.StringVar(&srv.Addr, "addr", ":8080", "listens on the TCP network address addr")
	flag.Uint64Var(&bulk, "bulk", 5e4, "count of numbers allocated by one request")
	flag.Uint64Var(&used, "used", 0, "last allocated number")
	flag.Parse()
	srv.Handler = handler{}

	done := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		if err := srv.Shutdown(context.Background()); err != nil {
			logger.Println("Shutdown error", err)
		}
		close(done)
	}()
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Println("ListenAndServe error", err)
	}
	<-done
}

type handler struct{}

func (handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	n := atomic.AddUint64(&used, 1)
	from, to := allocate(n, bulk)
	enc := json.NewEncoder(w)
	if err := enc.Encode(struct {
		From, To uint64
	}{from, to}); err != nil {
		logger.Println("Encode error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.Println("used", from, "to", to)
}

func allocate(n, bulk uint64) (from, to uint64) {
	if n == 0 {
		panic("n to be 0 is not allowed")
	}
	from = n*bulk - bulk
	to = n*bulk - 1
	return
}
