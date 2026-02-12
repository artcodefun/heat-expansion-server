package domain

import "time"

// NowUnixFunc returns current time as unix seconds. Overridable in tests.
var NowUnixFunc = func() int64 { return time.Now().Unix() }

// NowUnix returns current unix seconds via NowUnixFunc.
func NowUnix() int64 { return NowUnixFunc() }

// NowUnixNanoFunc returns current time as unix nanoseconds. Overridable in tests.
var NowUnixNanoFunc = func() int64 { return time.Now().UnixNano() }

// NowUnixNano returns current unix nanoseconds via NowUnixNanoFunc.
func NowUnixNano() int64 { return NowUnixNanoFunc() }
