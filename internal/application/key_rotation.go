// Package implementation for privacy transformation and sensitive-value protection.
package application

import (
	"context"
	"errors"
)

var ErrRotationNotStarted = errors.New("key rotation not started")

type KeyProvider interface {
	CurrentVersion(context.Context) (string, error)
	Rotate(context.Context) error
}
type TokenRevocation interface {
	Revoke(context.Context, string) error
}

type RotationResult struct {
	PreviousVersion string
	CurrentVersion  string
	Revoked         bool
}

type RotationCoordinator struct {
	keys        KeyProvider
	revocations TokenRevocation
}

func NewRotationCoordinator(keys KeyProvider, revocations TokenRevocation) *RotationCoordinator {
	return &RotationCoordinator{keys: keys, revocations: revocations}
}

func (c *RotationCoordinator) Rotate(ctx context.Context) (RotationResult, error) {
	previous, err := c.currentVersion(ctx)
	if err != nil {
		return RotationResult{}, err
	}
	if err := c.executeRotation(ctx); err != nil {
		return RotationResult{}, err
	}
	current, err := c.currentVersion(ctx)
	if err != nil {
		return RotationResult{}, err
	}
	if err := c.revokePreviousVersion(ctx, previous); err != nil {
		return RotationResult{}, err
	}
	return RotationResult{PreviousVersion: previous, CurrentVersion: current, Revoked: true}, nil
}
