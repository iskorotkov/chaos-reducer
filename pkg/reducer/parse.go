package reducer

import (
	"fmt"
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
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
			for _, step := range stage.Steps {
				template, ok := findTemplate(step.Name, manifest.Spec.Templates)
				if !ok {
					return Scenario{}, fmt.Errorf("couldn't find template for step")
				}

				phase, ok := findPhase(step.Name, manifest.Status.Nodes)
				if !ok {
					return Scenario{}, fmt.Errorf("couldn't find phase for step")
				}

				steps = append(steps, Step{
					Failed:       phase == v1alpha1.NodeFailed,
					WorkflowStep: step,
					Template:     template,
				})
			}

			stages = append(stages, Stage{steps})
		}
	}

	return Scenario{}, fmt.Errorf("couldn't find steps template")
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
