package zardoz

type RunContext struct {
	testFunction TestFunction
	contextName  string
	ranTest      Test
}

func (rc RunContext) printErrorHints() {
	for _, err := range rc.ranTest.Errors {
		err.printErrorHint()
	}
}

func (rc RunContext) printErrorLines() {
	for _, err := range rc.ranTest.Errors {
		err.printErrorLine()
	}
}
