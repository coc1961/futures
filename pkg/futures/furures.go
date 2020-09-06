package futures

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

type FutureFunction func(future FutureParam) (result interface{}, err error)

type Option func(*future)

func WithTimeout(timeout time.Duration) Option {
	return func(f *future) {
		f.timeout = timeout
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
}

type future struct {
	// Function
	futureFunc FutureFunction

	// Optional timeout
	timeout time.Duration

	// function return values
	value interface{}
	err   error

	// exit chan
	done chan bool

	// Context
	cancel context.CancelFunc
	ctx    context.Context
}

func New(fn FutureFunction, options ...Option) (Future, error) {
	if fn == nil {
		return nil, errors.New("function must not be null")
	}
	f := &future{
		futureFunc: fn,
		value:      nil,
		err:        nil,
		done:       make(chan bool, 10),
		timeout:    time.Duration(-1),
	}
	for _, op := range options {
		op(f)
	}
	return f.start(), nil
}

func (f *future) Done() <-chan struct{} {
	return f.ctx.Done()
}

func (f *future) Cancel() {
	if f.cancel != nil {
		f.cancel()
	}
}

func (f *future) Wait(d time.Duration) bool {
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

func (f *future) start() *future {
	if f.timeout == time.Duration(-1) {
		f.ctx, f.cancel = context.WithCancel(context.Background())
	} else {
		f.ctx, f.cancel = context.WithTimeout(context.Background(), f.timeout)
	}

	go f.run()

	return f
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

func (f *future) run() {
	defer func() {
		f.done <- true
	}()

	exitOk := make(chan bool)
	go func() {
		defer func() {
			if ret := recover(); ret != nil {
				f.err = f.error(ret)
				f.done <- true
			}
		}()
		f.value, f.err = f.futureFunc(f)
		exitOk <- true
	}()

	select {
	case <-exitOk:
	case <-f.ctx.Done():
		f.value = nil
		f.err = f.ctx.Err()
		return
	}
}
