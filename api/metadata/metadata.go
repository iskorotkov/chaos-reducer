package metadata

const (
	Prefix = "chaosframework.com"

	VersionV1 = "v1"

	TypeFailure = "failure"
	TypeUtility = "utility"

	SeverityNonCritical = "non critical"
	SeverityCritical    = "critical"
	SeverityLethal      = "lethal"

	ScaleContainer      = "container"
	ScalePod            = "pod"
	ScaleDeploymentPart = "deployment part"
	ScaleDeployment     = "deployment"
	ScaleNode           = "node"

	PhaseFalsePositiveCheck  = "false positive check"
	PhaseStageReduction      = "stage reduction"
	PhaseStageReductionCheck = "stage reduction check"
	PhaseStepReduction       = "step reduction"
	PhaseStepReductionCheck  = "step reduction check"
	PhaseFinished            = "finished"
)

type TemplateMetadata struct {
	Version  string `annotation:"version"`
	Type     string `annotation:"type"`
	Severity string `annotation:"severity"`
	Scale    string `annotation:"scale"`
}

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
