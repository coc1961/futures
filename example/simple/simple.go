package main

import (
	"context"
	"fmt"
	"time"

	"github.com/coc1961/futures/pkg/futures"
)

var t = Test{0}

func main() {
	th := make([]futures.Future, 0)
	for i := 0; i < 1000; i++ {
		f, _ := MakeFuture()
		th = append(th, f)
	}

	for i := 0; i < 100; i++ {
		th[i].Wait(time.Second * 20)
		r, _ := th[i].Result()
		fmt.Print(r, "-")
	}
}

func MakeFuture() (futures.Future, error) {
	return futures.New(context.Background(), t.run)
}

type Test struct {
	cont int
}

func (t *Test) run(future futures.FutureParam) (result interface{}, err error) {
	time.Sleep(5 * time.Second)
	t.cont++
	return t.cont, nil
}
