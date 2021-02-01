package zardoz

import "fmt"
import "github.com/fatih/color"

type TestFunction func(Test) Test

type Suite struct {
    runContexts []RunContext
}

type RunContext struct {
    testFunction TestFunction
    contextName string
}

func (s *Suite) Register(contextName string, testFunction TestFunction) {
    rc := RunContext { contextName: contextName, testFunction: testFunction }
    s.runContexts = append(s.runContexts, rc)
}

func runTest(rc RunContext) Test {
     t := Test {
        AssertCount: 0,
        Passes: 0,
        Failures: 0,
    }

    return rc.testFunction(t)
}

func (s Suite) Run() {
    ranTests := []Test{}

    fmt.Println("Running...")
    for _, rc := range s.runContexts {
        ranTest := runTest(rc)

        if ranTest.IsSuccessful() {
            color.Set(color.FgRed)
            fmt.Print("F")
        } else {
            color.Set(color.FgGreen)
            fmt.Print(".")
        }
        color.Unset()

        ranTests = append(ranTests, ranTest)
    }
    fmt.Println("\n... done!")

    for _, rt := range ranTests {
        rt.printResult()
    }

    fmt.Println()
    for _, rt := range ranTests {
        rt.printFailingLines()
    }
}
