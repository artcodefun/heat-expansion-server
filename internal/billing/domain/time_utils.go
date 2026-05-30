package domain

import "time"

// NowUnixFunc returns current time as unix seconds. Overridable in tests.
var NowUnixFunc = func() int64 { return time.Now().Unix() }

// NowUnix returns current unix seconds via NowUnixFunc.
func NowUnix() int64 { return NowUnixFunc() }
