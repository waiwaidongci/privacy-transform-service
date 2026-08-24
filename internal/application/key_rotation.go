// Package implementation for privacy transformation and sensitive-value protection.
package application

import "context"

type KeyProvider interface {
	CurrentVersion(context.Context) (string, error)
	Rotate(context.Context) error
}
type TokenRevocation interface {
	Revoke(context.Context, string) error
}
