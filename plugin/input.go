package plugin

import "context"

// Input is the interface that every data-collection plugin must implement.
//
// The SNMP Poller is the canonical first implementation. Future inputs (e.g.
// a gNMI subscriber or NETCONF listener) implement this same contract so the
// engine can manage them identically.
//
// Lifecycle:
//
//  1. The engine calls Start with a cancellable context and an output channel.
//  2. The plugin spawns its internal goroutines and begins sending Envelopes
//     on the output channel. It MUST respect context cancellation and stop
//     sending when the context is done.
//  3. The engine calls Stop for graceful shutdown. Stop blocks until all
//     in-flight work is drained and internal goroutines have exited.
//     After Stop returns the plugin MUST NOT send on the output channel.
//
// Thread safety: Start and Stop are called from the engine goroutine.
// The output channel may be shared among multiple Inputs, so sends must be
// non-blocking or respect context cancellation to avoid deadlocks.
type Input interface {
	// Name returns a short, unique, human-readable identifier for this input
	// plugin instance (e.g. "snmp_poller"). It is used in log messages,
	// metrics labels, and Envelope.Source.
	Name() string

	// Start begins data collection. Collected data is sent as Envelopes on
	// the provided output channel. The context controls the plugin's
	// lifetime — when it is cancelled the plugin must wind down promptly.
	//
	// Start must return quickly. Long-running work (polling loops, listeners)
	// should be launched in background goroutines. An error return means the
	// plugin failed to initialise and will not produce data.
	Start(ctx context.Context, out chan<- Envelope) error

	// Stop performs a graceful shutdown: completes in-flight operations,
	// flushes pending data, and releases resources. It blocks until
	// shutdown is complete. After Stop returns the plugin will not send
	// any further Envelopes.
	Stop()
}
