package zardoz

type describeFunction func(*Suite) 

func Describe(describe describeFunction) {
    suite:= Suite { runContexts: []RunContext{} }

    describe(&suite)
    suite.Run()
    
}
func NewSuite() Suite {
    return Suite { runContexts: []RunContext{} }
}
