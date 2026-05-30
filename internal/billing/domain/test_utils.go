package domain

import "testing"

// SetTestNow sets a deterministic time provider for tests and restores it on cleanup.
// It sets NowUnixFunc to return 'sec'.
func SetTestNow(t *testing.T, sec int64) {
	oldNow := NowUnixFunc
	NowUnixFunc = func() int64 { return sec }
	t.Cleanup(func() {
		NowUnixFunc = oldNow
	})
}
