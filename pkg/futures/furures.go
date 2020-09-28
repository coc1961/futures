package futures

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

type FutureFunction func(future FutureParam) (result interface{}, err error)

type Option func(*futureConfig)

func WithTimeout(timeout time.Duration) Option {
	return func(f *futureConfig) {
		f.timeout = &timeout
	}
}

type Future interface {
	Wait(d time.Duration) bool
	Result() (result interface{}, err error)
	Cancel()
}

type FutureParam interface {
	Cancel()
	Done() <-chan struct{}
	Context() context.Context
}

type future struct {
	// function return values
	value interface{}
	err   error

	// exit chan
	done chan bool

	// finished process?
	processExit bool

	// Context
	cancel context.CancelFunc
	ctx    context.Context
}

type futureConfig struct {
	// Optional timeout
	timeout *time.Duration
}

func New(ctx context.Context, fn FutureFunction, options ...Option) (Future, error) {
	if fn == nil {
		return nil, errors.New("function must not be null")
	}
	f := &future{
		value: nil,
		err:   nil,
		done:  make(chan bool, 10),
	}

	ff := &futureConfig{
		timeout: nil,
	}
	for _, op := range options {
		op(ff)
	}

	if ff.timeout == nil {
		f.ctx, f.cancel = context.WithCancel(ctx)
	} else {
		f.ctx, f.cancel = context.WithTimeout(ctx, *ff.timeout)
	}

	go f.run(fn)

	return f, nil
}

func (f *future) Done() <-chan struct{} {
	return f.ctx.Done()
}

func (f *future) Cancel() {
	if f.cancel != nil {
		f.cancel()
	}
}

func (f *future) Context() context.Context {
	return f.ctx
}

func (f *future) Wait(d time.Duration) bool {
	if f.processExit {
		return true
	}
	select {
	case <-f.done:
		return true
	case <-time.After(d):
		return false
	}

}
func (f *future) Result() (result interface{}, err error) {
	if f.value != nil || f.err != nil {
		return f.value, f.err
	}
	return nil, errors.New("running")
}

func (f *future) error(i interface{}) error {
	switch x := i.(type) {
	case error:
		return x
	case string:
		return errors.New(x)
	default:
		return errors.Errorf("%v", x)
	}
}

func (f *future) run(fn FutureFunction) {
	defer func() {
		f.done <- true
	}()

	exitOk := make(chan bool)
	go func(fn FutureFunction) {
		defer func() {
			if ret := recover(); ret != nil {
				f.err = f.error(ret)
				exitOk <- false
			}
		}()
		f.value, f.err = fn(f)
		exitOk <- true
	}(fn)

	select {
	case <-exitOk:
		f.processExit = true
	case <-f.ctx.Done():
		f.value = nil
		f.err = f.ctx.Err()
		return
	}
}
