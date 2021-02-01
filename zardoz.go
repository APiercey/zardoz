package zardoz

// import "fmt"
// import "github.com/fatih/color"


func NewSuite() Suite {
    return Suite { runContexts: []RunContext{} }
}

// func Run(testName string, testFunction TestFunction) {
//      t := Test {
//         AssertCount: 0,
//         Passes: 0,
//         Failures: 0,
//     }

//     color.Set(color.FgBlue)
//     fmt.Printf("Running %s...\n", testName)
//     color.Unset()
//     t = testFunction(t)
//     color.Set(color.FgBlue)
//     fmt.Printf("\nFinished running %s.\n", testName)
//     color.Unset()
//     t.printResult()

// }
