package futures

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFutures(t *testing.T) {

	/*
	******************************
	** Invalid Future Creation
	******************************
	 */
	_, err := New(nil)
	assert.NotNil(t, err)

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
	future, _ := New(fnFuture, WithTimeout(time.Second*10))

	// Wait for function end with 10 seconds timeout
	done := future.Wait(time.Second * 10)

	// Get Values Ok
	value, err := future.Result()

	assert.Equal(t, true, done)
	assert.Nil(t, err)
	assert.Equal(t, "Ok", value)

	/*
	******************************
	** Test 2 Timeout
	******************************
	 */

	// function that takes 10 seconds
	fnFuture = func(future FutureParam) (result interface{}, err error) {
		time.Sleep(time.Second * 10)
		return "Ok", nil
	}

	// Create a Future with timeout 2 seconds
	future, _ = New(fnFuture, WithTimeout(time.Second*2))

	// Wait for function end with 10 seconds timeout
	done = future.Wait(time.Second * 10)

	// Get Values Error
	value, err = future.Result()

	assert.Equal(t, true, done)
	assert.NotNil(t, err)
	assert.Equal(t, nil, value)

	/*
	******************************
	** Test 3 Function Cancel process
	******************************
	 */

	// function that takes 3 seconds and run cancel()
	fnFuture = func(future FutureParam) (result interface{}, err error) {
		time.Sleep(time.Second * 3)
		future.Cancel()
		time.Sleep(time.Second * 3)
		return "Ok", nil
	}

	// Create a Future with timeout 10 seconds
	future, _ = New(fnFuture, WithTimeout(time.Second*10))

	// Wait for function end with 10 seconds timeout
	done = future.Wait(time.Second * 10)

	// Get Values Error
	value, err = future.Result()

	assert.Equal(t, true, done)
	assert.NotNil(t, err)
	assert.Equal(t, nil, value)

	/*
	******************************
	** Test 4 Process Running
	******************************
	 */

	// function that takes 3 seconds and run cancel()
	fnFuture = func(future FutureParam) (result interface{}, err error) {
		time.Sleep(time.Second * 5)
		return "Ok", nil
	}

	// Create a Future with timeout 10 seconds
	future, _ = New(fnFuture, WithTimeout(time.Second*10))

	// Wait for function end with 1 seconds timeout
	done = future.Wait(time.Second * 1)

	// Get Values Error
	value, err = future.Result()

	assert.Equal(t, false, done)
	assert.NotNil(t, err)
	assert.Equal(t, nil, value)

	/*
	******************************
	** Test 5 Cancel Future
	******************************
	 */

	// function that takes 10 seconds and run cancel()
	fnFuture = func(future FutureParam) (result interface{}, err error) {
		time.Sleep(time.Second * 10)
		return "Ok", nil
	}

	// Create a Future with timeout 10 seconds
	future, _ = New(fnFuture, WithTimeout(time.Second*10))

	// Cancel Future
	future.Cancel()

	// Wait for function end with 1 seconds timeout
	done = future.Wait(time.Second * 10)

	// Get Values Error
	value, err = future.Result()

	assert.Equal(t, true, done)
	assert.NotNil(t, err)
	assert.Equal(t, nil, value)

}
