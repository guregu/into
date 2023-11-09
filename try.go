package into

import "fmt"

// Try executes fn, recovering and returning an error if it panics.
// If the panic value satisfies the error interface, it will be returned as-is;
// otherwise, it will be wrapped in [Panic].
func Try(fn func()) (err error) {
	defer catch(&err)
	fn()
	return
}

func catch(ep *error) {
	if ex := recover(); ex != nil {
		err, ok := ex.(error)
		if !ok {
			err = Panic{Value: ex}
		}
		*ep = err
	}
}

// Maybe tries to run the given function (such as [Int] or [String]),
// returning an error instead of panicking.
func Maybe[T any](try func(any, ...Option) T, value any, options ...Option) (result T, err error) {
	defer catch(&err)
	result = try(value, options...)
	return
}

// Panic is an error encapsulating a panicked value.
type Panic struct {
	// Value is the value passed to panic.
	Value any
}

// Error satisfies the error interface using fmt's %v verb to represent the panic value.
func (p Panic) Error() string {
	return fmt.Sprintf("panic: %v", p.Value)
}
