package environment

type Variable struct {
	Value string
}

func (v Variable) String() string {
	return "*****"
}

func New(value string) *Variable {
	return &Variable{value}
}
