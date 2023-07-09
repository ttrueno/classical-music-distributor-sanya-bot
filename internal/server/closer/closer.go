package closer

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/ttrueno/rl2-final/internal/lib/e"
)

var (
	ErrShutdownCanceled = errors.New("shutdown cancelled")
)

type closer struct {
	mu       sync.Mutex
	handlers []handler
}

func New() *closer {
	return &closer{}
}

func (c *closer) Add(h handler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlers = append(c.handlers, h)
}

func (c *closer) Close(ctx context.Context) (err error) {
	var (
		errmsg = `closer.Close`
	)
	defer func() { err = e.WrapIfErr(errmsg, err) }()

	c.mu.Lock()
	defer c.mu.Unlock()

	var (
		errs = make([]string, 0, len(c.handlers))
		done = make(chan struct{}, 1)
	)

	go func() {
		for _, h := range c.handlers {
			if err := h(ctx); err != nil {
				errs = append(
					errs,
					fmt.Sprintf("[!] execution stopped: %v\n", err),
				)
			}
		}

		done <- struct{}{}
	}()

	select {
	case <-done:
		break
	case <-ctx.Done():
		return ErrShutdownCanceled
	}

	if len(errs) != 0 {
		return fmt.Errorf(
			"shutdown finished with errors:\n%s",
			strings.Join(errs, "\n"),
		)
	}

	return nil
}

type handler func(ctx context.Context) error
