package main

import (
	"fmt"
	"time"

	"github.com/coc1961/futures/pkg/futures"
)

func main() {
	mf := NewMyFuture()

	for !mf.Done() {
		fmt.Println("Waiting...")
	}

	value, err := mf.Result()

	fmt.Println(value, err)
}

func NewMyFuture() *MyFuture {
	f := &MyFuture{}
	f.future, _ = futures.New(f.run)
	return f
}

type MyFuture struct {
	future futures.Future
}

func (m *MyFuture) run(future futures.FutureParam) (result interface{}, err error) {
	time.Sleep(time.Second * 10)
	return 1000, nil
}

func (m *MyFuture) Done() bool {
	return m.future.Wait(time.Second * 1)
}

func (m *MyFuture) Result() (int, error) {
	v, e := m.future.Result()
	return v.(int), e
}
