# into [![GoDoc](https://godoc.org/github.com/guregu/into?status.svg)](https://godoc.org/github.com/guregu/into)

**into** provides convenience functions for coercing dynamic values to concrete values.

Inlcudes:
- `into.String`, `into.Int`, `into.Uint`, `into.Float` for coercing `any` to their respective types
- `into.CanString`, `into.CanInt`, `into.CanUint`, `into.CanFloat` for testing coercibility
- `into.Try` and `into.Maybe` for catching panics and coercing them to `error`

### Motivation

This library is an experiment that might be useful if you have to deal with `map[string]any`, writing data mapper libraries, etc.
It's not intended to replace `strconv`; if you know your concrete types ahead of time just use the standard library.

The panicking nature of the API is an idea borrowed from everyone's favorite, the standard library `reflect` package.

## Example

```go
a := into.Int(42) // 42
b := into.Int("42", into.WithConvertString) // 42
c := into.Int(nil, into.WithFallback(1337)) // 1337

var ptr *int        // nil
d := into.Int(ptr)  // 0

ok := into.CanInt("blah") // false
err := into.Try(func() {
    a = into.Int("blah")
}) // into.ErrInvalid{Value: "blah", Type: "int"}
n, err := into.Maybe(into.Int, "blah") // 0, into.ErrInvalid{Value: "blah", Type: "int"}
```

## Performance

This library tries to avoid as much overhead as possible, and in many cases achieves zero allocation.

Benchmarks on a M2 Mac Mini:

```
BenchmarkFloat-12                	601534540	        1.997 ns/op	      0 B/op	      0 allocs/op
BenchmarkFloatWithOptions-12     	444859920	        2.692 ns/op	      0 B/op	      0 allocs/op
BenchmarkFloatFallback-12        	441497448	        2.699 ns/op	      0 B/op	      0 allocs/op
BenchmarkInt-12                  	572993503	        2.096 ns/op	      0 B/op	      0 allocs/op
BenchmarkIntWithOptions-12       	442683682	        2.715 ns/op	      0 B/op	      0 allocs/op
BenchmarkIntFallback-12          	440095513	        2.704 ns/op	      0 B/op	      0 allocs/op
BenchmarkString-12               	331075846	        3.579 ns/op	      0 B/op	      0 allocs/op
BenchmarkStringWithOptions-12    	307433682	        3.914 ns/op	      0 B/op	      0 allocs/op
BenchmarkStringFallback-12       	264048576	        4.490 ns/op	      0 B/op	      0 allocs/op
BenchmarkUint-12                 	621541999	        1.928 ns/op	      0 B/op	      0 allocs/op
BenchmarkUintWithOptions-12      	447733209	        2.673 ns/op	      0 B/op	      0 allocs/op
BenchmarkUintFallback-12         	445057358	        2.689 ns/op	      0 B/op	      0 allocs/op
```
