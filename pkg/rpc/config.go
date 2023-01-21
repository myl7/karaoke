// Copyright (C) 2023 myl7
// SPDX-License-Identifier: Apache-2.0

package rpc

// PConfig is per-server publicly available config
type PConfig struct {
	Addr  string
	PK    *[32]byte
	Layer int
}
