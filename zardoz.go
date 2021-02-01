package zardoz

import "fmt"
import "time"
import "runtime"
import "strings"
import "github.com/fatih/color"

type Test struct {
    AssertCount int
    Passes int
    Failures int
    Errors []string
}

type TestFunction func(Test) Test
type AssertFunction func() bool

func (t Test) printResult() {
    color.Red(strings.Join(t.Errors[:],"\n"))
    color.Blue("\n%d assertions ran with", t.AssertCount)
    color.Green("  %d passing", t.Passes)
    color.Red("  %d failing", t.Failures)
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
    code := readCode(fn, line - 1, line + 3)
    color.Red(errMessage)
    fmt.Println(code)

    t.addError(fmt.Sprintf("Assertion failed %s:%d", fn, line))
}

func (t *Test) addError(errMessage string) {
    t.Errors = append(t.Errors, errMessage)
}

func Run(testName string, testFunction TestFunction) {
     t := Test {
        AssertCount: 0,
        Passes: 0,
        Failures: 0,
    }

    color.Set(color.FgBlue)
    fmt.Printf("Running %s...\n", testName)
    color.Unset()
    t = testFunction(t)
    color.Set(color.FgBlue)
    fmt.Printf("\nFinished running %s.\n", testName)
    color.Unset()
    t.printResult()

}
