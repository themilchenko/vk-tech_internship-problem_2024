package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/themilchenko/vk-tech_internship-problem_2024/internal/app"
	"github.com/themilchenko/vk-tech_internship-problem_2024/internal/config"
)

func main() {
	// Parsing flags
	var configPath string
	config.ParseFlag(&configPath)
	flag.Parse()

	// Creating config instance
	cfg := config.NewConfig()
	if err := cfg.Open(configPath); err != nil {
		log.Printf("failed to open config file with path: %s", configPath)
	}

	doneCh := make(chan os.Signal, 1)
	signal.Notify(doneCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	s := app.NewServer(&http.Server{
		Addr:         cfg.Server.Address,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.IdleTimeout,
	}, cfg)
	go func() {
		if err := s.Start(); err != nil {
			log.Fatal("failed to start server")
		}
	}()

	<-doneCh
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Server.Shutdown(ctx); err != nil {
		log.Fatal("failed to stop server", err.Error())
		return
	}
}
