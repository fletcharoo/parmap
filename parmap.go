package parmap

import (
	"errors"
	"fmt"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type doFunc[IN, OUT any] func(input IN) (result OUT, err error)

type data[T any] struct {
	index int
	value T
	final bool
}

// ErrMap is a map[int]error where the key is the index of the input that failed.
// ErrMap implements Stringer.
// ErrMap implements Error.
type ErrMap map[int]error

// ErrJoin joins all the errors in the ErrMap in a deterministic way.
// Each error is of the form "<key>: <error>".
func (e ErrMap) ErrJoin() error {
	errs := make([]error, len(e))

	keys := maps.Keys(e)
	slices.Sort(keys)

	for i, k := range keys {
		errs[i] = fmt.Errorf("%d: %w", k, e[k])
	}

	return errors.Join(errs...)
}

// String joins all the errors in the ErrMap and returns it as a string.
func (e ErrMap) String() string {
	return e.ErrJoin().Error()
}

// Error joins all the errors in the ErrMap and returns it as a string.
func (e ErrMap) Error() string {
	return e.String()
}

func startResultActor[OUT any](doneChan chan struct{}, length int) (resultChan chan data[OUT], results []OUT) {
	results = make([]OUT, length)
	resultChan = make(chan data[OUT])

	go func() {
		for {
			select {
			case result, open := <-resultChan:
				if open {
					results[result.index] = result.value
					doneChan <- struct{}{}
				} else {
					return
				}
			}
		}
	}()

	return resultChan, results
}

func startErrActor(doneChan chan struct{}) (errChan chan data[error], erm ErrMap) {
	erm = make(ErrMap)
	errChan = make(chan data[error])

	go func() {
		for {
			select {
			case erd, open := <-errChan:
				if open {
					erm[erd.index] = erd.value
					doneChan <- struct{}{}
				} else {
					return
				}
			}
		}
	}()

	return errChan, erm
}

func startDoActor[IN, OUT any](do doFunc[IN, OUT], inputChan chan data[IN], resultChan chan data[OUT], errChan chan data[error]) {
	go func() {
		for {
			select {
			case input, open := <-inputChan:
				if open {
					result, err := do(input.value)

					if err != nil {
						errChan <- data[error]{
							index: input.index,
							value: err,
						}

						break
					}

					resultChan <- data[OUT]{
						index: input.index,
						value: result,
					}
				} else {
					return
				}
			}
		}
	}()
}

// Do runs the do func for each input in the inputs slice in parallel and
// returns a slice of results and an error.
// The length of the result slice will always be the same length as the inputs slice.
// If no errors occurred in the execution of the do funcs, the returned err will be nil.
// If errors occurred in the execution of the do funcs, the result slice will
// have the zero value of the OUT type at the indexes of the failed inputs.
func Do[IN, OUT any](inputs []IN, do doFunc[IN, OUT]) (result []OUT, err ErrMap) {
	inputsLen := len(inputs)
	inputChan := make(chan data[IN])
	doneChan := make(chan struct{})
	defer close(doneChan)
	resultChan, results := startResultActor[OUT](doneChan, inputsLen)
	defer close(resultChan)
	errChan, erm := startErrActor(doneChan)
	defer close(errChan)

	for range inputsLen {
		startDoActor(do, inputChan, resultChan, errChan)
	}

	for i, v := range inputs {
		inputChan <- data[IN]{
			index: i,
			value: v,
		}
	}

	close(inputChan)

	for range inputsLen {
		<-doneChan
	}

	if len(erm) != 0 {
		return results, erm
	}

	return results, nil
}
