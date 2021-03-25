package metadata

const (
	PhaseFalsePositiveCheck  = "false positive check"
	PhaseStageReduction      = "stage reduction"
	PhaseStageReductionCheck = "stage reduction check"
	PhaseStepReduction       = "step reduction"
	PhaseStepReductionCheck  = "step reduction check"
	PhaseFinished            = "finished"
)

type WorkflowMetadata struct {
	Version  string `annotation:"version"`
	Phase    string `annotation:"phase"`
	Original string `annotation:"original"`

	Iteration  int    `annotation:"iteration"`
	BasedOn    string `annotation:"based-on"`
	FollowedBy string `annotation:"followed-by"`

	Attempt  int    `annotation:"attempt"`
	Next     string `annotation:"next"`
	Previous string `annotation:"previous"`
}
