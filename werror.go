// Package werror provides a Result[T] type for elegant error handling in Go.
//
// It eliminates repetitive "if err != nil" checks through method chaining,
// while maintaining Go's explicit error handling philosophy.
//
// Example usage:
//
//	result := werror.Wrap(os.ReadFile("config.json")).
//	    ThenFn(parseJSON).
//	    ThenFn(validateConfig).
//	    OnError(func(err error) {
//	        log.Printf("failed: %v", err)
//	    })
//
//	if result.IsOk() {
//	    config := result.Value()
//	}
package werror

// Result represents a value that can be either successful (Ok) or failed (Err).
type Result[T any] struct {
	value T
	err   error
}

// Ok creates a successful Result containing the given value.
func Ok[T any](value T) Result[T] {
	return Result[T]{value: value}
}

// Err creates a failed Result containing the given error.
func Err[T any](err error) Result[T] {
	return Result[T]{err: err}
}

// Wrap converts a standard Go (value, error) pair into a Result[T].
// If err is not nil, returns Err. Otherwise returns Ok with the value.
func Wrap[T any](value T, err error) Result[T] {
	if err != nil {
		return Err[T](err)
	}

	return Ok(value)
}

// Unwrap extracts the underlying value and error from the Result.
func (r Result[T]) Unwrap() (T, error) {
	return r.value, r.err
}

// ThenFn chains a function that returns (T, error).
// If the Result is Err, the function is not called and the error propagates.
func (r Result[T]) ThenFn(fn func(T) (T, error)) Result[T] {
	if r.err != nil {
		return r
	}

	v, err := fn(r.value)
	return Wrap(v, err)
}

// Then chains a function that returns Result[T].
// If the Result is Err, the function is not called and the error propagates.
func (r Result[T]) Then(fn func(T) Result[T]) Result[T] {
	if r.err != nil {
		return r
	}

	return fn(r.value)
}

// OnSuccess executes the given function if the Result is Ok.
// Returns the Result unchanged for chaining.
func (r Result[T]) OnSuccess(fn func(T)) Result[T] {
	if r.err == nil {
		fn(r.value)
	}

	return r
}

// OnError executes the given function if the Result is Err.
// Returns the Result unchanged for chaining.
func (r Result[T]) OnError(fn func(error)) Result[T] {
	if r.err != nil {
		fn(r.err)
	}

	return r
}

// Always executes the given function regardless of the Result state.
// Returns the Result unchanged for chaining.
func (r Result[T]) Always(fn func()) Result[T] {
	fn()
	return r
}

// Value returns the underlying value (zero value if Result is Err).
func (r Result[T]) Value() T {
	return r.value
}

// Error returns the underlying error (nil if Result is Ok).
func (r Result[T]) Error() error {
	return r.err
}

// IsOk returns true if the Result is successful.
func (r Result[T]) IsOk() bool {
	return r.err == nil
}

// IsErr returns true if the Result contains an error.
func (r Result[T]) IsErr() bool {
	return r.err != nil
}

// Must returns the value if Ok, otherwise panics with the error.
// Use for cases where an error is genuinely unexpected.
func (r Result[T]) Must() T {
	if r.err != nil {
		panic(r.err)
	}

	return r.value
}

// Or returns the value if Ok, otherwise returns the provided fallback.
func (r Result[T]) Or(v T) T {
	if r.err != nil {
		return v
	}

	return r.value
}
