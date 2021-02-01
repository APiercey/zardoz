package zardoz
import "fmt"
import "github.com/fatih/color"

type Error struct {
    preview string
    errMessage string
    failingLine string
}

func (e Error) printErrorHint() {
    fmt.Println()
    color.Red("    %s", e.errMessage)
    fmt.Print(e.preview)
}

func (e Error) printErrorLine() {
    color.Red(e.failingLine)
}
