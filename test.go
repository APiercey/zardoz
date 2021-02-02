package zardoz

import "fmt"
import "time"
import "runtime"
import "sync"

type AssertFunction func() bool

type Test struct {
    AssertCount int
    Passes int
    Failures int
    AsyncAssertions sync.WaitGroup
    Errors []Error
}

func doLoopedTest(assertFunction AssertFunction, timeOutMilliseconds int) bool {
    result := false

    for start := time.Now(); time.Since(start) < time.Duration(timeOutMilliseconds) * time.Millisecond; {
        if assertFunction() { 
            result = true
            break
        } 
    }

    return result
}

func (t *Test) Assert(expectation bool) {
    t.AssertCount++

    if expectation {
        t.Passes++
    } else {
        t.addError(buildErrorFromCaller("Expected true when evaluating:"))
        t.Failures++
    }
}

func (t *Test) IsSuccessful() bool {
    return t.Failures == 0
}

func (t *Test) AssertSync(assertFunction AssertFunction, timeOutMilliseconds int) {
    t.AssertCount++

    if doLoopedTest(assertFunction, timeOutMilliseconds) {
        t.Passes++
    } else {
        errMessage := fmt.Sprintf("Never returned true after %dms when evaluating:", timeOutMilliseconds)
        t.addError(buildErrorFromCaller(errMessage))
        t.Failures++
    }
}


func (t *Test) AssertAsync(assertFunction AssertFunction, timeOutMilliseconds int) {
    t.AssertCount++

    preemptiveErr := buildErrorFromCaller(fmt.Sprintf("Never returned true after %dms when evaluating:", timeOutMilliseconds))

    t.AsyncAssertions.Add(1)
    go func(wg *sync.WaitGroup) {
        if doLoopedTest(assertFunction, timeOutMilliseconds) {
            t.Passes++
        } else {
            t.addError(preemptiveErr)
            t.Failures++
        }

        defer wg.Done()
    }(&t.AsyncAssertions)
}

func (t *Test) handleFailure(errMessage string) {
     _, fn, line, _ := runtime.Caller(2)
    code := readCode(fn, line, line + 3)

    err := Error { 
        preview: code, 
        errMessage: errMessage,
        failingLine: fmt.Sprintf("Assertion failed %s:%d", fn, line),
    }

    t.addError(err)
}

func (t *Test) addError(err Error) {
    t.Errors = append(t.Errors, err)
}

func buildErrorFromCaller(errMessage string) Error {
     _, fn, line, _ := runtime.Caller(2)
    code := readCode(fn, line, line + 3)

    return Error { 
        preview: code, 
        errMessage: errMessage,
        failingLine: fmt.Sprintf("Assertion failed %s:%d", fn, line),
    }
}
