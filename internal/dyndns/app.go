package dyndns

import (
	"time"

	"codeberg.org/maltech/desec-dyndns-client/internal/config"
)

type App struct {
	cfg *config.Config
}

func New(cfg *config.Config) *App {
	return &App{cfg: cfg}
}

func (a *App) Run() {
	ticker := time.NewTicker(a.cfg.Interval)
	defer ticker.Stop()

	a.runOnce()

	for range ticker.C {
		a.runOnce()
	}
}

func (a *App) runOnce() {
	a.update("A")
	a.update("AAAA")
}
