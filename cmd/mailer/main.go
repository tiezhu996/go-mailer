package main

import (
	"fmt"

	"mailer/internal/config"
	"mailer/internal/service"
	"mailer/internal/store"
)

func main() {
	cfg := config.Load()
	st := store.New()
	svc := service.New(st, cfg.BatchSize)
	_ = svc
	fmt.Println("mailer ready")
}
