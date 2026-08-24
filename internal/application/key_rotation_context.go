package application

import (
	"context"
	"fmt"
)

func (c *RotationCoordinator) currentVersion(_ context.Context) (string, error) {
	version, err := c.keys.CurrentVersion(context.Background())
	if err != nil {
		return "", fmt.Errorf("read current key version: %w", err)
	}
	return version, nil
}
