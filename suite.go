package zardoz

import "fmt"
import "github.com/fatih/color"
// import "time"

type TestFunction func(*Test) 

type Suite struct {
    description string
    runContexts []RunContext
}

func (s *Suite) Test(contextName string, testFunction TestFunction) {
    s.runContexts = append(s.runContexts, RunContext { contextName: contextName, testFunction: testFunction })
}

func printResultFromRanContexts(ranContexts []RunContext) {
    for idx := range ranContexts {
        rc := &ranContexts[idx]
        if rc.ranTest.IsSuccessful() { continue }

        color.Yellow("\n\nFailures:")
        fmt.Println()
        break
    }

    errorCount := 0
    successCount := 0
    for idx := range ranContexts {
        rc := &ranContexts[idx]

        if !rc.ranTest.IsSuccessful() {
            errorCount++

            color.Yellow(fmt.Sprintf("  %d) %s", errorCount, rc.contextName))
            rc.printErrorHints()
        } else {
            successCount++
        }
    }

    fmt.Println()
    for idx := range ranContexts {
        rc := &ranContexts[idx]
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

func printResultFromSuite(s *Suite) {
    for idx := range s.runContexts {
        rc := &s.runContexts[idx]
        if rc.ranTest.IsSuccessful() { continue }

        color.Yellow("\n\nFailures:")
        fmt.Println()
        break
    }

    errorCount := 0
    successCount := 0
    for idx := range s.runContexts {
        rc := &s.runContexts[idx]

        if !rc.ranTest.IsSuccessful() {
            errorCount++

            color.Yellow(fmt.Sprintf("  %d) %s %s", errorCount, s.description, rc.contextName))
            rc.printErrorHints()
        } else {
            successCount++
        }
    }

    fmt.Println()
    for idx := range s.runContexts {
        rc := &s.runContexts[idx]
        rc.printErrorLines()
    }

    fmt.Println()

    color.Set(color.FgBlue)
    fmt.Printf("%d tests ran with ", len(s.runContexts))
    
    color.Set(color.FgGreen)
    fmt.Printf("%d passes", successCount)

    color.Set(color.FgBlue)
    fmt.Print(" and ")

    color.Set(color.FgRed)
    fmt.Printf("%d failures.", errorCount)
}

func executeContext(rc *RunContext) {
     rc.ranTest = Test {
        AssertCount: 0,
        Passes: 0,
        Failures: 0,
    }

    rc.testFunction(&rc.ranTest)

    rc.ranTest.AsyncAssertions.Wait()
}

func (s Suite) Run() {
    fmt.Println("\nRunning...")

    for idx := range s.runContexts {
        rc := &s.runContexts[idx]
        executeContext(rc)

        if rc.ranTest.IsSuccessful() {
            color.Set(color.FgGreen)
            fmt.Print(".")
        } else {
            color.Set(color.FgRed)
            fmt.Print("F")
        }
        color.Unset()
    }

    printResultFromSuite(&s)
    fmt.Println()
}
