package futures

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFutures(t *testing.T) {

	// function that takes a second
	funcTime := time.Second * 1
	fnFuture := func(ctx context.Context, cancel context.CancelFunc) (result interface{}, err error) {
		time.Sleep(funcTime)
		return "Ok", nil
	}

	// Create a Future with timeout 2 seconds
	future := New(fnFuture, time.Second*2)

	// Wait 10 seconds for function end
	done := future.Wait(time.Second * 10)

	// Get Values Ok
	value, err := future.Values()

	assert.Equal(t, true, done)
	assert.Nil(t, err)
	assert.Equal(t, "Ok", value)

}
