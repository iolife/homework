package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func serverApp(ctx context.Context, stop <-chan struct{}) error {

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "okay")
		},
	))

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	go func() {
		<-stop
		s.Shutdown(ctx)

	}()

	return s.ListenAndServe()
}
func main() {
	c, cancel := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(c)
	_, stop := make(chan error), make(chan struct{})
	g.Go(func() error {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
		select {
		case sig := <-signalChannel:
			stop <- struct{}{}
			fmt.Printf("signal : %s\n", sig)
			cancel()
		case <-ctx.Done():
			stop <- struct{}{}
			fmt.Println("close signal ")
			return ctx.Err()
		}
		return nil
	})

	g.Go(func() error {
		return serverApp(ctx, stop)
	})

	if err := g.Wait(); err == nil {
		fmt.Println("Successfully fetched all URLs.")
	} else {
		log.Println(err)
	}

}
