package rpc

// RoundStartMsg is sent from coordinator to server.
// Use JSON to mashal it since Redis sub/pub transmit string.
type RoundStartMsg struct {
	Round int
}

// RoundEndMsg is sent from server to coordinator.
// For mashaling, see [RoundStartMsg].
type RoundEndMsg struct {
	Round int
	ID    string
}
