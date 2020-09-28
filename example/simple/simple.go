package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/coc1961/futures/pkg/futures"
)

const max = 300

func main() {
	var t = Test{0}

	// Make Futures
	th := make([]futures.Future, 0)
	for i := 0; i < max; i++ {
		f, _ := futures.New(context.Background(), t.run)
		th = append(th, f)
	}

	// Wait and get result
	i := 0
	for i < max {
		if !th[i].Wait(time.Microsecond) {
			continue
		}
		r, e := th[i].Result()
		if e != nil {
			fmt.Print(e, "-")
		} else {
			fmt.Print(r, "-")
		}
		i++
	}
	fmt.Println("\nExit")
}

type Test struct {
	cont int
}

func (t *Test) run(future futures.FutureParam) (result interface{}, err error) {
	time.Sleep(time.Duration(rand.ExpFloat64()*10) * time.Millisecond)
	t.cont++
	if t.cont == 100 {
		panic("(error 100)")
	}
	return t.cont, nil
}
