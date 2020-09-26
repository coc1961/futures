package futures

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {

	/*
	******************************
	** Invalid Future Creation
	******************************
	 */
	_, err := New(context.Background(), nil)
	assert.NotNil(t, err)
}

func TestFutures_Ok(t *testing.T) {

	/*
	******************************
	** Test 1 Run Ok
	******************************
	 */

	// function that takes a second
	fnFuture := func(future FutureParam) (result interface{}, err error) {
		time.Sleep(time.Second * 1)
		return "Ok", nil
	}

	// Create a Future with timeout 2 seconds
	future, _ := New(context.Background(), fnFuture, WithTimeout(time.Second*10))

	// Wait for function end with 10 seconds timeout
	done := future.Wait(time.Second * 10)

	// Get Values Ok
	value, err := future.Result()

	assert.Equal(t, true, done)
	assert.Nil(t, err)
	assert.Equal(t, "Ok", value)
}
func TestFutures_Timeout(t *testing.T) {

	/*
	******************************
	** Test 2 Timeout
	******************************
	 */

	// function that takes 10 seconds
	fnFuture := func(future FutureParam) (result interface{}, err error) {
		time.Sleep(time.Second * 10)
		return "Ok", nil
	}

	// Create a Future with timeout 2 seconds
	future, _ := New(context.Background(), fnFuture, WithTimeout(time.Second*2))

	// Wait for function end with 10 seconds timeout
	done := future.Wait(time.Second * 10)

	// Get Values Error
	value, err := future.Result()

	assert.Equal(t, true, done)
	assert.NotNil(t, err)
	assert.Equal(t, "context deadline exceeded", err.Error())
	assert.Equal(t, nil, value)
}
func TestFutures_CancelProcess(t *testing.T) {

	/*
	******************************
	** Test 3 Function Cancel process
	******************************
	 */

	// function that takes 3 seconds and run cancel()
	fnFuture := func(future FutureParam) (result interface{}, err error) {
		time.Sleep(time.Second * 3)
		future.Cancel()
		time.Sleep(time.Second * 3)
		return "Ok", nil
	}

	// Create a Future with timeout 10 seconds
	future, _ := New(context.Background(), fnFuture, WithTimeout(time.Second*10))

	// Wait for function end with 10 seconds timeout
	done := future.Wait(time.Second * 10)

	// Get Values Error
	value, err := future.Result()

	assert.Equal(t, true, done)
	assert.NotNil(t, err)
	assert.Equal(t, "context canceled", err.Error())
	assert.Equal(t, nil, value)
}
func TestFutures_Running(t *testing.T) {

	/*
	******************************
	** Test 4 Process Running
	******************************
	 */

	// function that takes 3 seconds and run cancel()
	fnFuture := func(future FutureParam) (result interface{}, err error) {
		time.Sleep(time.Second * 5)
		return "Ok", nil
	}

	// Create a Future with timeout 10 seconds
	future, _ := New(context.Background(), fnFuture, WithTimeout(time.Second*10))

	// Wait for function end with 1 seconds timeout
	done := future.Wait(time.Second * 1)

	// Get Values Error
	value, err := future.Result()

	assert.Equal(t, false, done)
	assert.NotNil(t, err)
	assert.Equal(t, "running", err.Error())
	assert.Equal(t, nil, value)
}
func TestFutures_Cancel(t *testing.T) {

	/*
	******************************
	** Test 5 Cancel Future
	******************************
	 */

	// function that takes 10 seconds and run cancel()
	fnFuture := func(future FutureParam) (result interface{}, err error) {
		time.Sleep(time.Second * 10)
		return "Ok", nil
	}

	// Create a Future with timeout 10 seconds
	future, _ := New(context.Background(), fnFuture, WithTimeout(time.Second*10))

	// Cancel Future
	future.Cancel()

	// Wait for function end with 1 seconds timeout
	done := future.Wait(time.Second * 10)

	// Get Values Error
	value, err := future.Result()

	assert.Equal(t, true, done)
	assert.NotNil(t, err)
	assert.Equal(t, "context canceled", err.Error())
	assert.Equal(t, nil, value)
}
func TestFutures_Panic(t *testing.T) {

	/*
	******************************
	** Test 6 Function Panic Recover
	******************************
	 */

	// function that takes 10 seconds and run cancel()
	fnFuture := func(future FutureParam) (result interface{}, err error) {
		time.Sleep(time.Second * 1)
		panic("function panic")
	}

	// Create a Future with timeout 10 seconds
	future, _ := New(context.Background(), fnFuture, WithTimeout(time.Second*10))

	// Wait for function end with 1 seconds timeout
	done := future.Wait(time.Second * 10)

	// Get Values Error
	value, err := future.Result()

	assert.Equal(t, true, done)
	assert.NotNil(t, err)
	assert.Equal(t, "function panic", err.Error())
	assert.Equal(t, nil, value)

}
