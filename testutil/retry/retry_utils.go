// Package retry contains utils methods .
package retry

import (
	"net"
	"time"
)

// DefaultRetryTimeout is the default timeout used for retry.
const DefaultRetryTimeout = time.Second

// AwaitForPort awaits for the port within the timeout.
func AwaitForPort(host, port string, timeout time.Duration) error {
	return WithTimeout(func() error {
		_, err := net.Dial("tcp", net.JoinHostPort(host, port))
		return err
	}, timeout)
}

// WithTimeout retries the operation within the timout.
func WithTimeout(operation func() error, timeout time.Duration) error {
	startTime := time.Now().UnixNano()
	for {
		err := operation()
		if err == nil {
			return nil
		}
		if time.Now().UnixNano()-startTime > timeout.Nanoseconds() {
			return err
		}
		time.Sleep(DefaultRetryTimeout)
	}
}
