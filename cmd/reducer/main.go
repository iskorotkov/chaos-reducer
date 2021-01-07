package main

import (
	"context"
	"fmt"
	"github.com/argoproj/argo/pkg/apiclient/workflow"
	argoApi "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	argoClient "github.com/argoproj/argo/pkg/client/clientset/versioned"
	"github.com/caarlos0/env"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type Params struct {
	ArgoServer  string `env:"ARGO_SERVER"`
	Development bool   `env:"DEVELOPMENT"`
}

func main() {
	var params Params
	err := env.Parse(&params)
	if err != nil {
		log.Fatal(err)
	}

	sugar, err := createLogger(params)
	if err != nil {
		log.Fatal(err)
	}

	defer syncLogger(sugar)

	listenToEvents(params, sugar)
}

func listenToEvents(params Params, logger *zap.SugaredLogger) {
	client, err := createArgoClient(params)
	if err != nil {
		logger.Fatal(err)
	}

	service, err := client.WatchWorkflows(context.Background(), &workflow.WatchWorkflowsRequest{})
	if err != nil {
		logger.Fatal(err)
	}

	clientset, err := createClientset()
	if err != nil {
		logger.Fatal(err)
	}

	for {
		event, err := service.Recv()
		if err != nil {
			logger.Fatal(err)
		}

		wf := event.Object
		status := wf.Status

		// Skip workflows that hasn't failed
		if status.Phase != "Failed" {
			continue
		}

		// Skip workflows that finished more than 60 seconds ago
		// if !status.FinishedAt.Add(time.Second * 60).After(time.Now()) {
		// 	continue
		// }

		fmt.Printf("%s: %s (%s - %s)\n", wf.Name, status.Phase, status.StartedAt.Format("02.01.2006 15:04:05 MST"), status.FinishedAt.Format("02.01.2006 15:04:05 MST"))

		_, ok := wf.Annotations["chaosframework.com/next"]
		if ok {
			// Workflow already has follow-up workflow
			continue
		}

		iteration, ok := wf.Annotations["chaosframework.com/iteration"]
		if !ok {
			// Workflow wasn't launched as a part of test
			continue
		}

		iterationNum, err := strconv.ParseInt(iteration, 10, 32)
		if err != nil {
			logger.Fatal(err)
		}

		wf.Annotations["chaosframework.com/iteration"] = strconv.FormatInt(iterationNum+1, 10)

		// TODO: launch new workflow

		// TODO: assign new workflow name
		name := "name"
		namespace := "namespace"

		wf.Annotations["chaosframework.com/next"] = fmt.Sprintf("%s/%s", namespace, name)

		patch := &argoApi.Workflow{
			ObjectMeta: v1.ObjectMeta{
				Annotations: wf.Annotations,
			},
		}

		data, err := json.Marshal(patch)
		if err != nil {
			logger.Fatal(err)
		}

		// Use MergePatchType because Argo uses it
		wf, err = clientset.ArgoprojV1alpha1().Workflows(wf.Namespace).Patch(wf.Name, types.MergePatchType, data)
		if err != nil {
			logger.Fatal(err)
		}
	}
}

func createArgoClient(params Params) (workflow.WorkflowServiceClient, error) {
	url := params.ArgoServer
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := workflow.NewWorkflowServiceClient(conn)
	return client, nil
}

func createClientset() (argoClient.Interface, error) {
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

	clientset, err := argoClient.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func createLogger(params Params) (*zap.SugaredLogger, error) {
	var logger *zap.Logger
	var err error
	if params.Development {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		return nil, err
	}

	sugar := logger.Sugar()
	return sugar, nil
}

func syncLogger(logger *zap.SugaredLogger) {
	err := logger.Sync()
	if err != nil {
		log.Fatal(err)
	}
}
