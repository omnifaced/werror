package werror

import (
	"errors"
	"fmt"
	"testing"
)

var (
	doubleFn = func(x int) (int, error) {
		return x * 2, nil
	}

	failFn = func(x int) (int, error) {
		return 0, errors.New("fail")
	}

	doubleResult = func(x int) Result[int] {
		return Ok(x * 2)
	}

	failResult = func(x int) Result[int] {
		return Err[int](errors.New("fail"))
	}
)

func assertOk(t *testing.T, result Result[int], expectedValue int) {
	t.Helper()
	if !result.IsOk() {
		t.Error("Expected Ok result")
	}
	if result.Value() != expectedValue {
		t.Errorf("Expected value %d, got %d", expectedValue, result.Value())
	}
}

func assertErr(t *testing.T, result Result[int]) {
	t.Helper()
	if !result.IsErr() {
		t.Error("Expected Err result")
	}
}

func assertErrorIs(t *testing.T, expected, actual error) {
	t.Helper()
	if !errors.Is(expected, actual) {
		t.Errorf("Expected error %v, got %v", expected, actual)
	}
}

func TestOk(t *testing.T) {
	result := Ok(42)

	if !result.IsOk() {
		t.Error("Expected Ok result")
	}

	if result.Value() != 42 {
		t.Errorf("Expected value 42, got %d", result.Value())
	}

	if result.Error() != nil {
		t.Error("Expected nil error")
	}
}

func TestErr(t *testing.T) {
	testErr := errors.New("test error")
	result := Err[int](testErr)

	if !result.IsErr() {
		t.Error("Expected Err result")
	}

	if !errors.Is(testErr, result.Error()) {
		t.Errorf("Expected error %v, got %v", testErr, result.Error())
	}

	if result.Value() != 0 {
		t.Errorf("Expected zero value, got %d", result.Value())
	}
}

func TestWrapOk(t *testing.T) {
	result := Wrap(42, nil)

	if !result.IsOk() {
		t.Error("Expected Ok result")
	}

	if result.Value() != 42 {
		t.Errorf("Expected value 42, got %d", result.Value())
	}
}

func TestWrapErr(t *testing.T) {
	testErr := errors.New("test error")
	result := Wrap(42, testErr)

	if !result.IsErr() {
		t.Error("Expected Err result")
	}

	if !errors.Is(testErr, result.Error()) {
		t.Errorf("Expected error %v, got %v", testErr, result.Error())
	}
}

func TestUnwrap(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		result := Ok(42)
		value, err := result.Unwrap()

		if err != nil {
			t.Errorf("Expected nil error, got %v", err)
		}

		if value != 42 {
			t.Errorf("Expected value 42, got %d", value)
		}
	})

	t.Run("Err", func(t *testing.T) {
		testErr := errors.New("test error")
		result := Err[int](testErr)
		value, err := result.Unwrap()

		if !errors.Is(err, testErr) {
			t.Errorf("Expected error %v, got %v", testErr, err)
		}

		if value != 0 {
			t.Errorf("Expected zero value, got %d", value)
		}
	})
}

func TestThenFn(t *testing.T) {
	t.Run("Ok chain", func(t *testing.T) {
		assertOk(t, Ok(10).ThenFn(doubleFn), 20)
	})

	t.Run("Err chain", func(t *testing.T) {
		assertErr(t, Ok(10).ThenFn(failFn))
	})

	t.Run("Err propagation", func(t *testing.T) {
		testErr := errors.New("original error")
		result := Err[int](testErr).ThenFn(doubleFn)
		assertErr(t, result)
		assertErrorIs(t, testErr, result.Error())
	})
}

func TestThen(t *testing.T) {
	t.Run("Ok chain", func(t *testing.T) {
		assertOk(t, Ok(10).Then(doubleResult), 20)
	})

	t.Run("Err chain", func(t *testing.T) {
		assertErr(t, Ok(10).Then(failResult))
	})

	t.Run("Err propagation", func(t *testing.T) {
		testErr := errors.New("original error")
		result := Err[int](testErr).Then(doubleResult)
		assertErr(t, result)
		assertErrorIs(t, testErr, result.Error())
	})
}

func TestOnSuccess(t *testing.T) {
	t.Run("Ok executes callback", func(t *testing.T) {
		executed := false
		Ok(42).OnSuccess(func(v int) {
			executed = true
		})

		if !executed {
			t.Error("Expected callback to be executed")
		}
	})

	t.Run("Err does not execute callback", func(t *testing.T) {
		executed := false
		Err[int](errors.New("test")).OnSuccess(func(v int) {
			executed = true
		})

		if executed {
			t.Error("Expected callback not to be executed")
		}
	})

	t.Run("Returns result for chaining", func(t *testing.T) {
		result := Ok(42).OnSuccess(func(v int) {})

		if !result.IsOk() {
			t.Error("Expected Ok result")
		}
	})
}

func TestOnError(t *testing.T) {
	t.Run("Err executes callback", func(t *testing.T) {
		executed := false
		Err[int](errors.New("test")).OnError(func(err error) {
			executed = true
		})

		if !executed {
			t.Error("Expected callback to be executed")
		}
	})

	t.Run("Ok does not execute callback", func(t *testing.T) {
		executed := false
		Ok(42).OnError(func(err error) {
			executed = true
		})

		if executed {
			t.Error("Expected callback not to be executed")
		}
	})

	t.Run("Returns result for chaining", func(t *testing.T) {
		result := Err[int](errors.New("test")).OnError(func(err error) {})

		if !result.IsErr() {
			t.Error("Expected Err result")
		}
	})
}

func TestAlways(t *testing.T) {
	t.Run("Ok executes callback", func(t *testing.T) {
		executed := false
		Ok(42).Always(func() {
			executed = true
		})

		if !executed {
			t.Error("Expected callback to be executed")
		}
	})

	t.Run("Err executes callback", func(t *testing.T) {
		executed := false
		Err[int](errors.New("test")).Always(func() {
			executed = true
		})

		if !executed {
			t.Error("Expected callback to be executed")
		}
	})
}

func TestIsOk(t *testing.T) {
	if !Ok(42).IsOk() {
		t.Error("Expected Ok result to return true")
	}

	if Err[int](errors.New("test")).IsOk() {
		t.Error("Expected Err result to return false")
	}
}

func TestIsErr(t *testing.T) {
	if !Err[int](errors.New("test")).IsErr() {
		t.Error("Expected Err result to return true")
	}

	if Ok(42).IsErr() {
		t.Error("Expected Ok result to return false")
	}
}

func TestValue(t *testing.T) {
	t.Run("Ok returns value", func(t *testing.T) {
		if Ok(42).Value() != 42 {
			t.Error("Expected value 42")
		}
	})

	t.Run("Err returns zero value", func(t *testing.T) {
		if Err[int](errors.New("test")).Value() != 0 {
			t.Error("Expected zero value")
		}
	})
}

func TestError(t *testing.T) {
	t.Run("Err returns error", func(t *testing.T) {
		testErr := errors.New("test error")
		if !errors.Is(testErr, Err[int](testErr).Error()) {
			t.Error("Expected test error")
		}
	})

	t.Run("Ok returns nil", func(t *testing.T) {
		if Ok(42).Error() != nil {
			t.Error("Expected nil error")
		}
	})
}

func TestMust(t *testing.T) {
	t.Run("Ok returns value", func(t *testing.T) {
		value := Ok(42).Must()

		if value != 42 {
			t.Errorf("Expected value 42, got %d", value)
		}
	})

	t.Run("Err panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic")
			}
		}()

		Err[int](errors.New("test")).Must()
	})
}

func TestOr(t *testing.T) {
	t.Run("Ok returns value", func(t *testing.T) {
		value := Ok(42).Or(99)

		if value != 42 {
			t.Errorf("Expected value 42, got %d", value)
		}
	})

	t.Run("Err returns fallback", func(t *testing.T) {
		value := Err[int](errors.New("test")).Or(99)

		if value != 99 {
			t.Errorf("Expected fallback 99, got %d", value)
		}
	})
}

func TestChaining(t *testing.T) {
	add10 := func(x int) (int, error) {
		return x + 10, nil
	}

	result := Ok(5).ThenFn(doubleFn).ThenFn(add10)
	assertOk(t, result, 20)
}

func TestComplexChaining(t *testing.T) {
	successCalled := false
	errorCalled := false
	alwaysCalled := false

	divide := func(a, b int) (int, error) {
		if b == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return a / b, nil
	}

	Wrap(divide(10, 2)).
		OnSuccess(func(v int) {
			successCalled = true
		}).
		OnError(func(err error) {
			errorCalled = true
		}).
		Always(func() {
			alwaysCalled = true
		})

	if !successCalled {
		t.Error("Expected success callback to be called")
	}

	if errorCalled {
		t.Error("Expected error callback not to be called")
	}

	if !alwaysCalled {
		t.Error("Expected always callback to be called")
	}
}

func TestErrorPropagation(t *testing.T) {
	originalErr := errors.New("original error")

	result := Err[int](originalErr).
		ThenFn(func(x int) (int, error) {
			return x * 2, nil
		}).
		Then(func(x int) Result[int] {
			return Ok(x + 1)
		})

	if !result.IsErr() {
		t.Error("Expected Err result")
	}

	if !errors.Is(originalErr, result.Error()) {
		t.Error("Expected original error to propagate")
	}
}
