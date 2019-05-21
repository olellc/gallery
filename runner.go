package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func RunUntilSignal(handler http.Handler, addr string) error {
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	sigs := []os.Signal{syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM}

	shutdownTimeout := time.Duration(0)

	return ListenAndServeUntilSignal(server, sigs, shutdownTimeout)
}

func ListenAndServeUntilSignal(server *http.Server, sigs []os.Signal,
	shutdownTimeout time.Duration) error {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, sigs...)
	defer signal.Stop(sigCh)

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		return err
	case <-sigCh:
	}

	ctx := context.Background()
	if shutdownTimeout > 0 {
		tmpCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		ctx = tmpCtx
		defer cancel()
	}

	// https://golang.org/pkg/net/http/#Server.Shutdown
	// "When Shutdown is called, Serve, ListenAndServe, and ListenAndServeTLS
	// immediately return ErrServerClosed."
	// There is also race condition between
	// server.ListenAndServe() and server.Shutdown(ctx)
	// which must be harmless since Go 1.11. Prior Go 1.11 receiving from
	// errCh may block forever. See
	// https://github.com/golang/go/issues/20239 for details.
	err1 := server.Shutdown(ctx)
	err2 := <-errCh

	if err1 != nil {
		return err1
	}

	if err2 == http.ErrServerClosed {
		return nil
	}

	return err2
}
