package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lacsar712/railblock/internal/app"
	"github.com/lacsar712/railblock/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	if err := cfg.Validate(); err != nil {
		log.Fatalf("config: %v", err)
	}

	state := app.NewState(cfg)
	fmt.Fprint(os.Stderr, app.StartupBanner(state))

	if err := app.Run(state); err != nil {
		log.Fatalf("railblock: %v", err)
	}
}
