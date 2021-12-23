package testutil

import (
	"net"
	"time"
)

// DefaultRetryTimeout is the default timeout used for retry.
const DefaultRetryTimeout = time.Second

// AwaitForPort awaits for the port within the timeout.
func AwaitForPort(host, port string, timeout time.Duration) error {
	return RetryWithTimeout(func() error {
		_, err := net.Dial("tcp", net.JoinHostPort(host, port))
		return err
	}, timeout)
}

// RetryWithTimeout retries the operation within the timout.
func RetryWithTimeout(operation func() error, timeout time.Duration) error {
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
