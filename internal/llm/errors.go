package llm

import "errors"

// ErrRecoverableInference marks transport failures that occurred after Quine
// assembled a provider request but before a durable assistant response existed.
// Retrying from the same provider context is safe for Quine state.
var ErrRecoverableInference = errors.New("recoverable inference transport failure")

type RecoverableInferenceError struct {
	Err error
}

func (e *RecoverableInferenceError) Error() string {
	if e == nil || e.Err == nil {
		return ErrRecoverableInference.Error()
	}
	return e.Err.Error()
}

func (e *RecoverableInferenceError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func (e *RecoverableInferenceError) Is(target error) bool {
	return target == ErrRecoverableInference
}

func newRecoverableInferenceError(err error) error {
	if err == nil {
		return ErrRecoverableInference
	}
	if errors.Is(err, ErrRecoverableInference) {
		return err
	}
	return &RecoverableInferenceError{Err: err}
}
