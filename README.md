# werror

🚀 Elegant error handling for Go through Result[T] type and method chaining.

Eliminate repetitive `if err != nil` boilerplate while maintaining Go's explicit error handling philosophy.

## 📦 Installation

**Requirements:** Go 1.18+ (generics support required)

```bash
go get github.com/omnifaced/werror
```

## ⚡ Quick Start

```go
import "github.com/omnifaced/werror"

result := werror.Wrap(os.ReadFile("config.json")).
    ThenFn(parseJSON).
    ThenFn(validateConfig).
    OnError(func(err error) {
        log.Printf("failed to load config: %v", err)
    })

if result.IsOk() {
    config := result.Value()
    // use config
}
```

## 💡 Why werror?

**Before:**
```go
file, err := os.ReadFile("config.json")
if err != nil {
    return err
}

data, err := parseJSON(file)
if err != nil {
    return err
}

config, err := validateConfig(data)
if err != nil {
    return err
}
```

**After:**
```go
result := werror.Wrap(os.ReadFile("config.json")).
    ThenFn(parseJSON).
    ThenFn(validateConfig)

config, err := result.Unwrap()
```

## 📚 API

### Constructors

- **`Ok[T](value T)`** - create successful Result
- **`Err[T](err error)`** - create failed Result
- **`Wrap[T](value T, err error)`** - convert `(value, error)` pair to Result

### Chaining

- **`ThenFn(fn func(T) (T, error))`** - chain standard Go functions
- **`Then(fn func(T) Result[T])`** - chain Result-returning functions

### Side Effects

- **`OnSuccess(fn func(T))`** - execute on success
- **`OnError(fn func(error))`** - execute on error
- **`Always(fn func())`** - execute regardless of state

### Inspection

- **`IsOk()`** - check if successful
- **`IsErr()`** - check if failed
- **`Value()`** - get value (zero if error)
- **`Error()`** - get error (nil if ok)
- **`Unwrap()`** - get both `(value, error)`

### Utilities

- **`Must()`** - unwrap or panic
- **`Or(fallback T)`** - unwrap or return fallback

## 🔥 Examples

### Basic Usage

```go
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}

result := werror.Wrap(divide(10, 2))

if result.IsOk() {
    fmt.Println("Result:", result.Value())
}
```

### Chaining with ThenFn

```go
result := werror.Ok(10).
    ThenFn(double).
    ThenFn(square).
    ThenFn(validate)

value, err := result.Unwrap()
```

### Chaining with Then

```go
func processValue(x int) werror.Result[int] {
    if x < 0 {
        return werror.Err[int](fmt.Errorf("negative value"))
    }
    return werror.Ok(x * 2)
}

result := werror.Ok(10).
    Then(processValue).
    Then(validateRange)
```

### Side Effects

```go
werror.Wrap(processFile("data.txt")).
    OnSuccess(func(data string) {
        log.Println("Processing succeeded")
    }).
    OnError(func(err error) {
        log.Printf("Error: %v", err)
    }).
    Always(func() {
        log.Println("Operation completed")
    })
```

### Fallback Values

```go
config := werror.Wrap(loadConfig()).Or(defaultConfig)

data := werror.Wrap(fetchData()).Must()
```

## 🎯 Design Philosophy

- **Zero runtime overhead** - simple wrapper around value and error
- **Explicit, not magical** - errors don't disappear, just handled elegantly
- **Composable** - all methods return Result for chaining
- **Idiomatic Go** - works with standard `(value, error)` pattern

## ⚖️ License

Copyright (c) 2025 Abbasov Ravan (omnifaced)
All rights reserved.
