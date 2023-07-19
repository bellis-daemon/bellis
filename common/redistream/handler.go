package redistream

import "context"

type MessageHandler func(ctx context.Context, message *Message) error
