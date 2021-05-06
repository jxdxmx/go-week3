package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() (err error) {
		var once = sync.Once{}
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("hello go week3"))
		})
		server := &http.Server{
			Addr:              ":8080",
			Handler:           mux,
			TLSConfig:         nil,
			ReadTimeout:       0,
			ReadHeaderTimeout: 0,
			WriteTimeout:      0,
			IdleTimeout:       0,
			MaxHeaderBytes:    0,
			TLSNextProto:      nil,
			ConnState:         nil,
			ErrorLog:          nil,
			BaseContext:       nil,
			ConnContext:       nil,
		}
		defer func() { _ = server.Shutdown(ctx) }()
		for {
			select {
			case <-ctx.Done():
				return errors.New("ctx done ..111")
			default:
				go once.Do(func() {
					err = server.ListenAndServe()
				})
				if err != nil {
					return err
				}
				time.Sleep(time.Second)
			}
		}
	})
	g.Go(func() (err error) {
		var once = sync.Once{}
		c := make(chan os.Signal)
		for {
			select {
			case <-ctx.Done():
				return errors.New("ctx done ..222")
			default:
				once.Do(func() {
					for {
						//signal.Notify(c, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT)
						s := <-c
						txt := s.String()
						if txt == "interrupt" || txt == "quit" || txt == "aborted" ||
							txt == "killed" || txt == "terminated" {
							fmt.Println("sig:", txt)
							err = errors.New("sig end")
							break
						}
					}
				})
				if err != nil {
					return err
				}
				time.Sleep(time.Second)
			}
		}
	})
	err := g.Wait()
	fmt.Println("err group end, result:", err)
}
