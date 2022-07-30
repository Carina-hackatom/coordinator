package types

import (
	"context"
)

// Client describes the interface of Tendermint RPC client implementations.
type Client interface {
	Start() error
	Stop() error
	IsRunning() bool
	EventsClient
}

// EventsClient is reactive, you can subscribe to any message, given the proper string.
type EventsClient interface {
	Subscribe(ctx context.Context, query string) error
	Unsubscribe(ctx context.Context, query string) error
	// UnsubscribeAll unsubscribes from all the queries.
	UnsubscribeAll(ctx context.Context) error
}
