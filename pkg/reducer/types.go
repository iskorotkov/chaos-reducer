package reducer

import (
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
)

type Step struct {
	Failed       bool
	WorkflowStep v1alpha1.WorkflowStep
	Template     v1alpha1.Template
}

type Stage struct {
	Steps []Step
}

type Scenario struct {
	Stages   []Stage
	Workflow v1alpha1.Workflow
}
