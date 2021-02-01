package zardoz

import "fmt"
import "github.com/fatih/color"
// import "strings"
import "time"
import "runtime"

type AssertFunction func() bool

type Test struct {
    AssertCount int
    Passes int
    Failures int
    Errors []Error
}

func (t Test) printResult() {
    for _, err := range t.Errors {
        fmt.Println()
        color.Red(err.errMessage)
        fmt.Print(err.preview)
    }
}

func (t Test) printFailingLines() {
    for _, err := range t.Errors {
        color.Red(err.failingLine)
    }
}

func (t *Test) Assert(expectation bool) {
    t.AssertCount++

    if expectation {
        t.Passes++
    } else {
        t.handleFailure("Expected true when evaluating:");
        t.Failures++
    }
}

func (t Test) IsSuccessful() bool {
    return t.Failures > 0
}

func (t *Test) AssertAsync(assertFunction AssertFunction, timeOutMilliseconds int) {
    t.AssertCount++

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

    if result {
        t.Passes++
    } else {
        t.handleFailure(fmt.Sprintf("Never returned true after %dms when evaluating:", timeOutMilliseconds));
        t.Failures++
    }
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
