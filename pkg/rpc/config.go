package rpc

// PConfig is per-server publicly available config
type PConfig struct {
	Addr  string
	PK    *[32]byte
	Layer int
}
