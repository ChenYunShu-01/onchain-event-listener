package main

import (
	"flag"
	"os"

	"github.com/reddio-com/red-adapter/config"
	"github.com/reddio-com/red-adapter/core/executor"
)

var (
	cfgFile string
)

func init() {
	flag.StringVar(&cfgFile, "config", "config/adapter.toml", "config file")
}
func main() {
	cfg, err := config.LoadConfig(cfgFile)

	if err != nil {
		panic(err)
	}

	dsn := os.Getenv("DSN")

	cfg.DSN = dsn

	s, err := executor.NewExecutor(cfg)

	if err != nil {
		panic(err)
	}

	s.StartToWatchEvent()
}
