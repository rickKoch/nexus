package signals

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
)

func Context() context.Context {
	ctx, cancel := context.WithCancelCause(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		cancel(errors.New("cancelling context, received signal " + sig.String()))
		sig = <-sigCh
	}()

	return ctx
}
