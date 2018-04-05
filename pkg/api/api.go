package api

import (
	"os"
	"os/signal"
	"syscall"

	`github.com/spacelavr/pandora/pkg/log`
)

// Daemon start api daemon
func Daemon() bool {

	log.Debug("start api daemon")

	var (
		sig  = make(chan os.Signal, 1)
		done = make(chan bool, 1)
	)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-sig:
				done <- true
				return
			}
		}
	}()

	<-done
	log.Debug("handle SIGINT and SIGTERM")
	return true
}