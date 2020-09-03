package futures

import (
	"context"
	"errors"
	"time"
)

type FutureFunction func(ctx context.Context, cancel context.CancelFunc) (result interface{}, err error)

type Future interface {
	Wait(d time.Duration) bool
	Values() (result interface{}, err error)
}

type future struct {
	fn    FutureFunction
	value interface{}
	err   error
	done  bool
}

func New(fn FutureFunction, timeout time.Duration) Future {
	f := &future{
		fn:   fn,
		done: false,
	}
	return f.start(timeout)
}

func (f *future) Wait(d time.Duration) bool {
	runFunc := func() chan bool {
		ret := make(chan bool)
		go func() {
			for {
				if f.done {
					ret <- true
				}
				time.Sleep(time.Millisecond * 100)
			}
		}()
		return ret
	}

	select {
	case <-runFunc():
		return true
	case <-time.After(d):
		return false
	}

}
func (f *future) Values() (result interface{}, err error) {
	if f.done {
		return f.value, f.err
	}
	return nil, errors.New("running")
}

func (f *future) start(timeout time.Duration) *future {
	value := make(chan interface{})
	err := make(chan error)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	_ = cancel

	go func(fn FutureFunction, ctx context.Context, cancel context.CancelFunc, value chan<- interface{}, err chan<- error) {
		runFunc := func() chan bool {
			ret := make(chan bool)
			go func() {
				v, e := fn(ctx, cancel)
				value <- v
				err <- e
				ret <- true
			}()
			return ret
		}
		select {
		case <-runFunc():
		case <-ctx.Done():
			f.err = ctx.Err()
			return
		}
	}(f.fn, ctx, cancel, value, err)

	go func(f *future, ctx context.Context, value <-chan interface{}, err <-chan error) {
		defer func() {
			f.done = true
		}()

		select {
		case f.value = <-value:
			e := <-err
			if f.err == nil {
				f.err = e
			}
			return
		case <-ctx.Done():
			f.err = ctx.Err()
			return
		}
	}(f, ctx, value, err)

	return f
}
