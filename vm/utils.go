package vm

import (
	"context"
)

func isDone(ctx context.Context) (done bool) {
	select {
	case <-ctx.Done():
		return true

	default:
		return false
	}
}

func boolToByte(val bool) byte {
	if val {
		return 1
	}

	return 0
}
