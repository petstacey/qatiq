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

func (a *api) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", a.cfg.port),
		Handler:      a.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second, // TODO: remove magic
		WriteTimeout: 30 * time.Second, // TODO: remove magic
	}
	shutdownError := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		q := <-quit
		a.logRequest(a.service.Name, fmt.Sprintf("caught signal %s", map[string]string{
			"signal": q.String()}))
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}
		a.wg.Wait()
		shutdownError <- nil
	}()
	a.logRequest(a.service.Name, fmt.Sprintf("starting %s server", a.cfg.env))
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	err = <-shutdownError
	if err != nil {
		return err
	}
	a.logRequest(a.service.Name, fmt.Sprintf("%s server stopped", a.service.Name))
	return nil
}
