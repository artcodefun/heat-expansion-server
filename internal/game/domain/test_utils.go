package domain

import "testing"

// SetTestNow sets deterministic time providers for tests and restores them on cleanup.
// It sets NowUnixFunc to return 'sec' and NowUnixNanoFunc to return sec*1e9 + 123.
func SetTestNow(t *testing.T, sec int64) {
	oldNow := NowUnixFunc
	oldNowNano := NowUnixNanoFunc
	NowUnixFunc = func() int64 { return sec }
	NowUnixNanoFunc = func() int64 { return sec*1_000_000_000 + 123 }
	t.Cleanup(func() {
		NowUnixFunc = oldNow
		NowUnixNanoFunc = oldNowNano
	})
}
