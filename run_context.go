package zardoz

import "fmt"
import "github.com/fatih/color"

type RunContext struct {
    testFunction TestFunction
    contextName string
    ranTest Test
}

func (rc RunContext) printErrorHints() {
    for _, err := range rc.ranTest.Errors {
        fmt.Println()
        color.Red(err.errMessage)
        fmt.Print(err.preview)
    }
}

func (rc RunContext) printErrorLines() {
    for _, err := range rc.ranTest.Errors {
        color.Red(err.failingLine)
    }
}

