package promhttputil

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// DeferredRoundTripper is a wrapper around http.RoundTripper that
// lazily initializes the RoundTripper using a factory function.
type DeferredRoundTripper struct {
	factory func() (http.RoundTripper, error)

	loaded atomic.Bool
	mu     sync.RWMutex

	rt  http.RoundTripper
	err error
}

func (d *DeferredRoundTripper) load() {
	if d.loaded.Load() {
		return
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if d.loaded.Load() {
		return
	}

	defer d.loaded.Store(true)

	rt, err := d.factory()
	if err != nil {
		d.rt, d.err = nil, fmt.Errorf("error creating round tripper: %w", err)
		time.AfterFunc(5*time.Second, func() {
			d.mu.Lock()
			defer d.mu.Unlock()

			d.loaded.Store(false)
			d.rt, d.err = nil, nil
		})
		return
	}
	d.rt, d.err = rt, nil
}

// RoundTrip implements the http.RoundTripper interface.
func (d *DeferredRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	d.load()

	if d.err != nil {
		return nil, d.err
	}

	return d.rt.RoundTrip(req)
}

// NewDeferredRoundTripper creates a new DeferredRoundTripper with the given
// factory function. The factory function will be called the first time
// RoundTrip is called, and the resulting RoundTripper will be used for all
// subsequent calls.
func NewDeferredRoundTripper(factory func() (http.RoundTripper, error)) *DeferredRoundTripper {
	return &DeferredRoundTripper{
		factory: factory,
	}
}
