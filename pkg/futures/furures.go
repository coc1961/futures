package futures

import (
	"context"
	"errors"
	"time"
)

type FutureFunction func(future FutureParam) (result interface{}, err error)

type OptionFN func(*future)

func WithTimeout(timeout time.Duration) OptionFN {
	return func(f *future) {
		f.timeout = timeout
	}
}

type Future interface {
	FutureParam
	Wait(d time.Duration) bool
	Result() (result interface{}, err error)
}

type FutureParam interface {
	Cancel()
	Context() *context.Context
}

type future struct {
	// Function
	fn FutureFunction

	// Optional timeout
	timeout time.Duration

	// function return values
	value interface{}
	err   error

	// function return values chan
	valueCh chan interface{}
	errCh   chan error

	// exit chan
	done chan bool

	// Context
	cancel context.CancelFunc
	ctx    context.Context
}

func New(fn FutureFunction, options ...OptionFN) (Future, error) {
	if fn == nil {
		return nil, errors.New("function must not be null")
	}
	f := &future{
		fn:      fn,
		value:   nil,
		err:     nil,
		done:    make(chan bool, 1),
		valueCh: make(chan interface{}, 1),
		errCh:   make(chan error, 1),
		timeout: time.Duration(-1),
	}
	for _, op := range options {
		op(f)
	}

	return f.start(), nil
}

func (f *future) Context() *context.Context {
	return &f.ctx
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

	go f.waitFunctionEnd()
	go f.testExitConditions()

	return f
}

func (f *future) waitFunctionEnd() {
	ret := make(chan bool)
	go func() {
		v, e := f.fn(f)
		f.valueCh <- v
		f.errCh <- e
		ret <- true
	}()

	select {
	case <-ret:
	case <-f.ctx.Done():
		f.value = nil
		f.err = f.ctx.Err()
		return
	}
}

func (f *future) testExitConditions() {
	defer func() {
		f.done <- true
	}()

	select {
	case f.value = <-f.valueCh:
		e := <-f.errCh
		if f.err == nil {
			f.err = e
		}
		if f.err != nil {
			f.value = nil
		}
		return
	case <-f.ctx.Done():
		f.value = nil
		f.err = f.ctx.Err()
		return
	}

}
