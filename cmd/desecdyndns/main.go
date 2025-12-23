package main

import (
	"codeberg.org/maltech/desec-dyndns-client/internal/config"
	"codeberg.org/maltech/desec-dyndns-client/internal/dyndns"
	"codeberg.org/maltech/desec-dyndns-client/internal/metrics"
)

func main() {
	cfg := config.Load()

	metrics.Init()
	metrics.Start(cfg.MetricsAddr)

	app := dyndns.New(cfg)
	app.Run()
}
