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

func doAsyncTest(assertFunction AssertFunction, timeOutMilliseconds int) bool {
    result := false
    loopCount := 0

    for range time.Tick(10 * time.Millisecond) {
        if loopCount > (timeOutMilliseconds / 10) {
            break
        }
        if assertFunction() { 
            result = true
            break
        } 

        loopCount++
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

    if doAsyncTest(assertFunction, timeOutMilliseconds) {
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
        if doAsyncTest(assertFunction, timeOutMilliseconds) {
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
