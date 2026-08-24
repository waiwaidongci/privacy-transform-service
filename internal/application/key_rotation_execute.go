package application

import (
	"context"
	"fmt"
)

func (c *RotationCoordinator) executeRotation(_ context.Context) error {
	if err := c.keys.Rotate(context.Background()); err != nil {
		return fmt.Errorf("rotate key provider: %w", err)
	}
	return nil
}
