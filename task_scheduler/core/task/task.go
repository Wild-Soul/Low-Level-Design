package task

import (
	"context"
	"time"
)

type Task interface {
	Execute(context.Context) (int, time.Duration, error)
	GetId() int
}
