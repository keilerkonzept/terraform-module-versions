package util

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

// DoWithRetry runs the specified action. If it returns a value, return that value. If it returns an error, sleep for
// sleepBetweenRetries and try again, up to a maximum of maxRetries retries. If maxRetries is exceeded, return a
// MaxRetriesExceeded error.
func DoWithRetry(actionDescription string, maxRetries int, sleepBetweenRetries time.Duration, logger *logrus.Entry, action func() error) error {
	for i := 0; i <= maxRetries; i++ {
		logger.Printf(actionDescription)

		err := action()
		if err == nil {
			return nil
		}

		logger.Errorf("%s returned an error: %s. Sleeping for %s and will try again.", actionDescription, err.Error(), sleepBetweenRetries)
		time.Sleep(sleepBetweenRetries)
	}

	return MaxRetriesExceeded{Description: actionDescription, MaxRetries: maxRetries}
}

// MaxRetriesExceeded is an error that occurs when the maximum amount of retries is exceeded.
type MaxRetriesExceeded struct {
	Description string
	MaxRetries  int
}

func (err MaxRetriesExceeded) Error() string {
	return fmt.Sprintf("'%s' unsuccessful after %d retries", err.Description, err.MaxRetries)
}
