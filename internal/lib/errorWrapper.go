package lib

import "fmt"

func ErrWrapper(err error, op string) error {
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
