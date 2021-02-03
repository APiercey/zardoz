package zardoz

type describeFunction func(*Suite)

func Describe(description string, describe describeFunction) {
	suite := Suite{
		runContexts:    []RunContext{},
		description:    description,
		beforeFunction: func() {},
		afterFunction:  func() {},
	}

	describe(&suite)
	suite.Run()
}

func NewSuite() Suite {
	return Suite{runContexts: []RunContext{}}
}
