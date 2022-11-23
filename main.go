package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/tyjet/webhook-test/internal/config"
	"github.com/tyjet/webhook-test/internal/server"
)

func main() {
	cfg, err := config.Load("config", []string{"/etc/webhook-test"}, config.YAML)
	if err != nil {
		fmt.Printf("failed to load config error: %s", err.Error())
		os.Exit(1)
	}
	srv := server.NewServer(cfg)
	srv.Register()
	srv.Start()

	wait()
}

func wait() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
