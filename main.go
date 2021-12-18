package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func HelloServer(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}


func main() {
	g, ctx := errgroup.WithContext(context.Background())

	srv := &http.Server{
		Addr: ":8080",
	}
	g.Go(func() error {
		http.HandleFunc("/hello", HelloServer)
		fmt.Println("http server start")
		err := srv.ListenAndServe()
		return err
	})

	g.Go(func() error {
		<- ctx.Done()
		fmt.Println("http server stop")
		return srv.Shutdown(ctx)
	})

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("signal ctx done")
				return ctx.Err()
			case sig := <-signals:
				return errors.New("\nget signal " + sig.String() + ", application will shutdown\n")
			}
		}
	})

	if err := g.Wait(); err != nil {
		fmt.Println(err)
	}

	fmt.Println(ctx.Err())
}
