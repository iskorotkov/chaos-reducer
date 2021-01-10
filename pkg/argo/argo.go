package argo

import (
	"context"
	"encoding/json"
	"github.com/argoproj/argo/pkg/apiclient/workflow"
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/argoproj/argo/pkg/client/clientset/versioned"
	"google.golang.org/grpc"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

type Workflow *v1alpha1.Workflow

type Client struct {
	argoClient       workflow.WorkflowServiceClient
	kubernetesClient versioned.Interface
}

func NewClient(argoHost string) (*Client, error) {
	argo, err := createArgoClient(argoHost)
	if err != nil {
		return nil, err
	}

	kubernetes, err := createClientset()
	if err != nil {
		return nil, err
	}

	return &Client{
		argoClient:       argo,
		kubernetesClient: kubernetes,
	}, nil
}

func (c *Client) CreateWorkflow(ctx context.Context, request *workflow.WorkflowCreateRequest) (Workflow, error) {
	return c.argoClient.CreateWorkflow(ctx, request)
}

func (c *Client) GetWorkflow(ctx context.Context, request *workflow.WorkflowGetRequest) (Workflow, error) {
	return c.argoClient.GetWorkflow(ctx, request)
}

func (c *Client) SetAnnotations(name, namespace string, annotations map[string]string) (Workflow, error) {
	patch := &v1alpha1.Workflow{
		ObjectMeta: v1.ObjectMeta{
			Annotations: annotations,
		},
	}

	data, err := json.Marshal(patch)
	if err != nil {
		return nil, err
	}

	// Use MergePatchType because Argo uses it
	wf, err := c.kubernetesClient.ArgoprojV1alpha1().Workflows(namespace).Patch(name, types.MergePatchType, data)
	if err != nil {
		return nil, err
	}

	return wf, nil
}

func (c *Client) SetLabels(name, namespace string, labels map[string]string) (Workflow, error) {
	patch := &v1alpha1.Workflow{
		ObjectMeta: v1.ObjectMeta{
			Labels: labels,
		},
	}

	data, err := json.Marshal(patch)
	if err != nil {
		return nil, err
	}

	// Use MergePatchType because Argo uses it
	wf, err := c.kubernetesClient.ArgoprojV1alpha1().Workflows(namespace).Patch(name, types.MergePatchType, data)
	if err != nil {
		return nil, err
	}

	return wf, nil
}

func (c *Client) WatchWorkflows(ctx context.Context, request *workflow.WatchWorkflowsRequest) (workflow.WorkflowService_WatchWorkflowsClient, error) {
	return c.argoClient.WatchWorkflows(ctx, request)
}

func createArgoClient(target string) (workflow.WorkflowServiceClient, error) {
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := workflow.NewWorkflowServiceClient(conn)
	return client, nil
}

func createClientset() (versioned.Interface, error) {
	var (
		config *rest.Config
		err    error
	)

	if os.Getenv("KUBERNETES_HOST") != "" {
		config, err = rest.InClusterConfig()
	} else {
		configFile := filepath.Join(homedir.HomeDir(), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", configFile)
	}

	if err != nil {
		return nil, err
	}

	clientset, err := versioned.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}
