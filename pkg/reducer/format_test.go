package reducer

import (
	"fmt"
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/iskorotkov/chaos-reducer/api/metadata"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	"testing"
)

func TestFormat(t *testing.T) {
	type args struct {
		workflow v1alpha1.Workflow
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "empty workflow",
			args: args{workflow: v1alpha1.Workflow{
				ObjectMeta: v1.ObjectMeta{
					Name:      "name",
					Namespace: "namespace",
					Labels:    map[string]string{},
					Annotations: map[string]string{
						prefixed("version"):     metadata.VersionV1,
						prefixed("phase"):       metadata.PhaseStepReduction,
						prefixed("original"):    "abc",
						prefixed("iteration"):   "8",
						prefixed("based-on"):    "def",
						prefixed("followed-by"): "ghi",
						prefixed("attempt"):     "10",
						prefixed("next"):        "klm",
						prefixed("previous"):    "opq",
					},
				},
				Spec: v1alpha1.WorkflowSpec{
					Entrypoint: "steps-template-123",
					Templates: []v1alpha1.Template{
						{Name: "steps-template-123", Steps: []v1alpha1.ParallelSteps{}},
					},
				},
			}},
		},
		{
			name: "workflow with steps",
			args: args{workflow: v1alpha1.Workflow{
				ObjectMeta: v1.ObjectMeta{
					Name:      "name",
					Namespace: "namespace",
					Labels:    map[string]string{},
					Annotations: map[string]string{
						prefixed("version"):     metadata.VersionV1,
						prefixed("phase"):       metadata.PhaseStageReduction,
						prefixed("original"):    "some original",
						prefixed("iteration"):   "8",
						prefixed("based-on"):    "some ref",
						prefixed("followed-by"): "some other ref",
						prefixed("attempt"):     "10",
						prefixed("next"):        "some next ref",
						prefixed("previous"):    "some previous ref",
					},
				},
				Spec: v1alpha1.WorkflowSpec{
					Entrypoint: "steps-template-123",
					Templates: []v1alpha1.Template{
						{
							Name: "steps-template-123",
							Steps: []v1alpha1.ParallelSteps{
								{
									Steps: []v1alpha1.WorkflowStep{
										{Name: "step-1-1", Template: "step-1-1"},
										{Name: "step-1-2", Template: "step-1-2"},
									},
								},
								{
									Steps: []v1alpha1.WorkflowStep{
										{Name: "step-2-1", Template: "step-2-1"},
										{Name: "step-2-2", Template: "step-2-2"},
									},
								},
							},
						},
						{
							Name: "step-1-1",
							Metadata: v1alpha1.Metadata{
								Labels: map[string]string{},
								Annotations: map[string]string{
									prefixed("version"):  metadata.VersionV1,
									prefixed("type"):     metadata.TypeFailure,
									prefixed("severity"): metadata.SeverityNonCritical,
									prefixed("scale"):    metadata.ScaleContainer,
								},
							},
						},
						{
							Name: "step-1-2",
							Metadata: v1alpha1.Metadata{
								Labels: map[string]string{},
								Annotations: map[string]string{
									prefixed("version"):  metadata.VersionV1,
									prefixed("type"):     metadata.TypeFailure,
									prefixed("severity"): metadata.SeverityCritical,
									prefixed("scale"):    metadata.ScalePod,
								},
							},
						},
						{
							Name: "step-2-1",
							Metadata: v1alpha1.Metadata{
								Labels: map[string]string{},
								Annotations: map[string]string{
									prefixed("version"):  metadata.VersionV1,
									prefixed("type"):     metadata.TypeFailure,
									prefixed("severity"): metadata.SeverityLethal,
									prefixed("scale"):    metadata.ScaleDeployment,
								},
							},
						},
						{
							Name: "step-2-2",
							Metadata: v1alpha1.Metadata{
								Labels: map[string]string{},
								Annotations: map[string]string{
									prefixed("version"):  metadata.VersionV1,
									prefixed("type"):     metadata.TypeFailure,
									prefixed("severity"): metadata.SeverityNonCritical,
									prefixed("scale"):    metadata.ScaleNode,
								},
							},
						},
					},
				},
				Status: v1alpha1.WorkflowStatus{
					Nodes: map[string]v1alpha1.NodeStatus{
						"step-1-1": {Name: "step-1-1", Phase: v1alpha1.NodeSucceeded},
						"step-1-2": {Name: "step-1-2", Phase: v1alpha1.NodeSucceeded},
						"step-2-1": {Name: "step-2-1", Phase: v1alpha1.NodeSucceeded},
						"step-2-2": {Name: "step-2-2", Phase: v1alpha1.NodeSucceeded},
					},
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scenario, err := Parse(&tt.args.workflow)
			if err != nil {
				t.Fatalf("error occured while parsing: %s", err)
			}

			workflow, err := Format(scenario)
			if err != nil {
				t.Fatalf("error occurred while formatting: %s", err)
			}

			got, want := *workflow, tt.args.workflow
			want.Name = ""
			want.Status = v1alpha1.WorkflowStatus{}

			if got.Name != "" {
				t.Errorf("new workflow has non-empty name")
			}

			if !reflect.DeepEqual(got.Status, v1alpha1.WorkflowStatus{}) {
				t.Errorf("new workflow has non-empty status")
			}

			if !reflect.DeepEqual(got, want) {
				t.Errorf("Format(Parse()) got = %#v,\n want %#v", got, want)
			}
		})
	}
}

func prefixed(key string) string {
	return fmt.Sprintf("%s/%s", metadata.Prefix, key)
}
