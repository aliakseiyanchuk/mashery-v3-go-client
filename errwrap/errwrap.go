package errwrap

import "fmt"

type WrappedError struct {
	Context string
	Cause   error
}

func (w *WrappedError) Error() string {
	return fmt.Sprintf("%s: %s", w.Context, w.Cause)
}

func (m *WrappedError) Unwrap() error {
	return m.Cause
}
