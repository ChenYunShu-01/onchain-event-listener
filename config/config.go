package config

import "github.com/BurntSushi/toml"

type Chain struct {
	RPC     string `toml:"rpc"`
	ChainID int    `toml:"chain_id"`
}

type Cfg struct {
	DSN           string            `toml:"dsn"`
	Chains        map[string]*Chain `toml:"chains"`
	Chain         string            `toml:"chain"`
	EventBlockGap uint64            `toml:"event_block_gap"`
}

func NewConfig() *Cfg {
	return &Cfg{}
}

func LoadConfig(filepath string) (*Cfg, error) {
	cfg := &Cfg{}
	_, err := toml.DecodeFile(filepath, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
