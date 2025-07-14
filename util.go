package task_chain

import (
	"github.com/pkg/errors"
)

func GetParameter[T any](t TaskInterface, key string) (T, error) {
	value := t.getParameter(key)
	if value == nil {
		return *new(T), errors.Errorf("get parameter error: can not found key - %s", key)
	}

	v, ok := value.(T)
	if !ok {
		return *new(T), errors.Errorf("get parameter error: type error (key: %s)", key)
	}

	return v, nil
}
