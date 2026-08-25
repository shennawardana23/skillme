package server

import (
	"fmt"

	"fixture/pkg/config"
)

// Start begins serving requests using cfg. Logs the effective timeout on
// startup so operators can confirm the deployed configuration.
func Start(cfg *config.Config) {
	fmt.Printf("starting server on %s:%d with timeout %s\n", cfg.Host, cfg.Port, cfg.Timeout.String())
}
