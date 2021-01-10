package reducer

import (
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
)

func Format(scenario Scenario) (*v1alpha1.Workflow, error) {
	res := scenario.Workflow.DeepCopy()

	templates := []v1alpha1.Template{{
		Name:  "entrypoint",
		Steps: buildStepsTemplate(scenario),
	}}

	for _, stage := range scenario.Stages {
		for _, step := range stage.Steps {
			for _, candidate := range scenario.Workflow.Spec.Templates {
				if candidate.Name != step.Template.Name {
					continue
				}

				templates = append(templates, candidate)
			}
		}
	}

	res.Spec.Templates = templates
	return res, nil
}

func buildStepsTemplate(scenario Scenario) []v1alpha1.ParallelSteps {
	stages := make([]v1alpha1.ParallelSteps, 0)
	for _, stage := range scenario.Stages {
		steps := make([]v1alpha1.WorkflowStep, 0)
		for _, step := range stage.Steps {
			steps = append(steps, step.WorkflowStep)
		}

		stages = append(stages, v1alpha1.ParallelSteps{Steps: steps})
	}

	return stages
}
