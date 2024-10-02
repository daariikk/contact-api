package e

import "fmt"

func Err(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}
