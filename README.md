# futures

Future, Golang Implementation

> IMPORTANT example project !!

Example:

```go

// Future function that takes 10 seconds
fnFuture := func(future FutureParam) (result interface{}, err error) {
    time.Sleep(time.Second * 10)
    return "Ok", nil
}

// Create a Future with timeout 20 seconds
future, err := New(fnFuture, WithTimeout(time.Second*20))

// Wait for function end with 10 seconds timeout
done := future.Wait(time.Second * 10)

// Get Values Error
value, err := future.Result()

```
