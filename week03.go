package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

func serverApp(ctx context.Context, stop <-chan struct{}) error {
	h := func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Hello, world!\n")
	}

	s := http.Server{Addr: "0.0.0.0:8080", Handler: h}
	go func() {
		<-stop
		s.Shutdown(ctx)

	}()
	return s.ListenAndServe()
}
func main() {

	g, _ := errgroup.WithContext(context.Background())
	done, stop := make(chan error), make(chan struct{})
	g.Go(func() error {
		return serverApp(stop)
	})
	if err := g.Wait(); err == nil {
		fmt.Println("Successfully fetched all URLs.")
	}
	go func() {
		time.Sleep(time.Second * 5)
		fmt.Println("Context canceled")

	}()

}
