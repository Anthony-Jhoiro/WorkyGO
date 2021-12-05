package configParser

type ExecutionEnvironment struct {
	Env    map[string]string
	Params map[string]string
}

func LoadExecutionEnvironment(env map[string]string, params map[string]string) *ExecutionEnvironment {
	environment := &ExecutionEnvironment{}

	environment.Env = env

	environment.Params = params

	return environment
}
