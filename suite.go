package zardoz

import "fmt"
import "github.com/fatih/color"

type TestFunction func(*Test) 

type Suite struct {
    runContexts []RunContext
}

func (s *Suite) Register(contextName string, testFunction TestFunction) {
    rc := RunContext { contextName: contextName, testFunction: testFunction }
    s.runContexts = append(s.runContexts, rc)
}

func printResultFromRanContexts(ranContexts []RunContext) {
    for _, rc := range ranContexts {
        if rc.ranTest.IsSuccessful() { continue }

        color.Yellow("\n\nFailures:")
        fmt.Println()
        break
    }

    errorCount := 0
    successCount := 0
    for _, rc := range ranContexts {
        if !rc.ranTest.IsSuccessful() {
            errorCount++

            color.Yellow(fmt.Sprintf("  %d) %s", errorCount, rc.contextName))
            rc.printErrorHints()
        } else {
            successCount++
        }
    }

    fmt.Println()
    for _, rc := range ranContexts {
        rc.printErrorLines()
    }

    fmt.Println()

    color.Set(color.FgBlue)
    fmt.Printf("%d tests ran with ", len(ranContexts))
    
    color.Set(color.FgGreen)
    fmt.Printf("%d passes", successCount)

    color.Set(color.FgBlue)
    fmt.Print(" and ")

    color.Set(color.FgRed)
    fmt.Printf("%d failures.", errorCount)


}

func executeContext(rc RunContext) RunContext {
     rc.ranTest = Test {
        AssertCount: 0,
        Passes: 0,
        Failures: 0,
    }

    rc.testFunction(&rc.ranTest)

    for _, c := range rc.ranTest.AsyncAssertions {
        <- c
    }

    return rc
}

func (s Suite) Run() {
    ranContexts := []RunContext{}

    fmt.Println("\nRunning...")
    for _, rc := range s.runContexts {
        ranContext := executeContext(rc)

        if ranContext.ranTest.IsSuccessful() {
            color.Set(color.FgGreen)
            fmt.Print(".")
        } else {
            color.Set(color.FgRed)
            fmt.Print("F")
        }
        color.Unset()

        ranContexts = append(ranContexts, ranContext)
    }

    printResultFromRanContexts(ranContexts)
    fmt.Println()
}
