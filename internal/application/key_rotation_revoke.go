package application

import (
	"context"
	"fmt"
)

func (c *RotationCoordinator) revokePreviousVersion(_ context.Context, version string) error {
	if err := c.revocations.Revoke(context.Background(), version); err != nil {
		return fmt.Errorf("revoke previous key version %q: %w", version, err)
	}
	return nil
}
