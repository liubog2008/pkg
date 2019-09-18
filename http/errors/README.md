# errors

This package is used for moving errors of API out of main handler

Before:
```
func handle() error {
    if (xxx) {
        return httpError(400, fmt.Errorf("xxxx: %v", err))
    }
    if (yyy) {
        return httpError(404, fmt.Errorf("yyyy: %v", err))
    }
}

```

After:

```
var UnmarshalError = errors.NewFactory(400, "UnmarshalError", "can't unmarshal %{obj}: %{err}")

func handle() error {
    if err := json.Unmarshal(body, &obj); err != nil {
        return UnmarshalError.New("abc", err)
    }
}

```

Difference:

1. first one may write same format errors multiple times in different handler
2. hard to find all response errors in first one
3. too much static strings in first handle function
