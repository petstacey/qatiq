package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (b *broker) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", b.cfg.port),
		Handler:      b.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second, // TODO: remove magic
		WriteTimeout: 30 * time.Second, // TODO: remove magic
	}
	shutdownError := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		q := <-quit
		fmt.Printf("caught signal %s", map[string]string{
			"signal": q.String(),
		})
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}
		// app.logger.PrintInfo("completing background tasks", map[string]string{
		// 	"addr": srv.Addr,
		// })
		b.wg.Wait()
		shutdownError <- nil
	}()
	// s.logger.PrintInfo("starting server", map[string]string{
	// 	"addr": srv.Addr,
	// 	"env":  s.cfg.Env,
	// })
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	err = <-shutdownError
	if err != nil {
		return err
	}
	// app.logger.PrintInfo("server stopped", map[string]string{
	// 	"addr": srv.Addr,
	// })
	return nil
}
