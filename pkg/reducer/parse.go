package reducer

import (
	"errors"
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	api "github.com/iskorotkov/chaos-reducer/api/metadata"
	"github.com/iskorotkov/chaos-reducer/pkg/metadata"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	ErrStepTemplate        = errors.New("couldn't find template for step")
	ErrStepsTemplate       = errors.New("couldn't find steps template")
	ErrMetadata            = errors.New("couldn't marshall or unmarshall metadata")
	ErrUnknownTemplateType = errors.New("template type is unknown")
)

func Parse(manifest *v1alpha1.Workflow) (Scenario, error) {
	templates := manifest.Spec.Templates
	for _, t := range templates {
		if t.GetType() != v1alpha1.TemplateTypeSteps {
			continue
		}

		stages := make([]Stage, 0)
		for _, stage := range t.Steps {
			steps := make([]Step, 0)
			utilitySteps := make([]UtilityStep, 0)

			for _, step := range stage.Steps {
				template, ok := findTemplate(step.Name, manifest.Spec.Templates)
				if !ok {
					return Scenario{}, ErrStepTemplate
				}

				phase, ok := findPhase(step.Name, manifest.Status.Nodes)
				if !ok {
					phase = NodeNotReached
				}

				objectMeta := v1.ObjectMeta{
					Labels:      template.Metadata.Labels,
					Annotations: template.Metadata.Annotations,
				}

				var templateMeta api.TemplateMetadata
				err := metadata.Unmarshal(objectMeta, &templateMeta, api.Prefix)
				if err != nil {
					return Scenario{}, ErrMetadata
				}

				switch templateMeta.Type {
				case api.TypeFailure:
					steps = append(steps, NewStep(phase, step, template, templateMeta))
				case api.TypeUtility:
					utilitySteps = append(utilitySteps, NewUtilityStep(step, template, templateMeta))
				default:
					return Scenario{}, ErrUnknownTemplateType
				}
			}

			stages = append(stages, NewStage(steps, utilitySteps))
		}

		var workflowMeta api.WorkflowMetadata
		err := metadata.Unmarshal(manifest.ObjectMeta, &workflowMeta, api.Prefix)
		if err != nil {
			return Scenario{}, ErrMetadata
		}

		return NewScenario(stages, *manifest, workflowMeta), nil
	}

	return Scenario{}, ErrStepsTemplate
}

func findTemplate(name string, templates []v1alpha1.Template) (v1alpha1.Template, bool) {
	for _, t := range templates {
		if t.Name == name {
			return t, true
		}
	}

	return v1alpha1.Template{}, false
}

func findPhase(name string, nodesStatus v1alpha1.Nodes) (v1alpha1.NodePhase, bool) {
	for _, n := range nodesStatus {
		if n.Name == name {
			return n.Phase, true
		}
	}

	return "", false
}
