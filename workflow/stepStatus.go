package workflow

const (
	StepOK      StepStatus = 0
	StepFail               = 1
	StepRunning            = 2
	StepPending            = 3
)

func (s StepStatus) GetName() string {
	switch s {
	case StepOK:
		return "OK"
	case StepPending:
		return "Pending"
	case StepFail:
		return "Fail"
	case StepRunning:
		return "Running..."
	default:
		return "Unknown state"
	}
}
