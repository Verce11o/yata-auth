package email

import (
	"context"
)

type EmailPublisher interface {
	Publish(ctx context.Context, message []byte) error
}
