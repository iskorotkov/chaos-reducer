package reducer

import (
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	api "github.com/iskorotkov/chaos-reducer/api/metadata"
	"github.com/iskorotkov/chaos-reducer/pkg/metadata"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Format(scenario Scenario) (*v1alpha1.Workflow, error) {
	res := scenario.Workflow.DeepCopy()
	res.ObjectMeta.Name = ""
	res.Status = v1alpha1.WorkflowStatus{}

	err := metadata.Marshal(&res.ObjectMeta, &scenario.Metadata, api.Prefix)
	if err != nil {
		return nil, ErrMetadata
	}

	templates := []v1alpha1.Template{{
		Name:  scenario.Workflow.Spec.Entrypoint,
		Steps: buildStepsTemplate(scenario),
	}}

	for _, stage := range scenario.Stages {
		for _, step := range stage.Steps {
			for _, candidate := range scenario.Workflow.Spec.Templates {
				if candidate.Name != step.Template.Name {
					continue
				}

				var objectMeta v1.ObjectMeta
				err := metadata.Marshal(&objectMeta, &step.Metadata, api.Prefix)
				if err != nil {
					return nil, ErrMetadata
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
