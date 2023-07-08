package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	game := New()

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGKILL)
		sig := <-ch
		game.Shutdown(sig.String())
	}()

	game.Run()
}
