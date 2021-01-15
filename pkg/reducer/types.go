package reducer

import (
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/iskorotkov/chaos-reducer/api/metadata"
)

const (
	NodeNotReached = v1alpha1.NodePhase("NotReached")
)

type Step struct {
	Phase        v1alpha1.NodePhase
	WorkflowStep v1alpha1.WorkflowStep
	Template     v1alpha1.Template
	Metadata     metadata.TemplateMetadata
}

func NewStep(p v1alpha1.NodePhase, ws v1alpha1.WorkflowStep, t v1alpha1.Template, m metadata.TemplateMetadata) Step {
	return Step{Phase: p, WorkflowStep: ws, Template: t, Metadata: m}
}

func (s Step) Failed() bool {
	return s.Phase == v1alpha1.NodeFailed
}

type UtilityStep struct {
	WorkflowStep v1alpha1.WorkflowStep
	Template     v1alpha1.Template
	Metadata     metadata.TemplateMetadata
}

func NewUtilityStep(ws v1alpha1.WorkflowStep, t v1alpha1.Template, m metadata.TemplateMetadata) UtilityStep {
	return UtilityStep{WorkflowStep: ws, Template: t, Metadata: m}
}

type Stage struct {
	Steps        []Step
	UtilitySteps []UtilityStep
}

func NewStage(steps []Step, utilitySteps []UtilityStep) Stage {
	return Stage{Steps: steps, UtilitySteps: utilitySteps}
}

type Scenario struct {
	Stages   []Stage
	Workflow v1alpha1.Workflow
	Metadata metadata.WorkflowMetadata
}

func NewScenario(stages []Stage, workflow v1alpha1.Workflow, metadata metadata.WorkflowMetadata) Scenario {
	return Scenario{Stages: stages, Workflow: workflow, Metadata: metadata}
}
