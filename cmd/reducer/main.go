package main

import (
	"context"
	"fmt"
	"github.com/argoproj/argo/pkg/apiclient/workflow"
	"github.com/caarlos0/env"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
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

	service := startService(params, sugar.Named("startup"))
	listenToEvents(service, sugar.Named("loop"))
}

func listenToEvents(service workflow.WorkflowService_WatchWorkflowsClient, sugar *zap.SugaredLogger) {
	for {
		event, err := service.Recv()
		if err != nil {
			sugar.Fatal(err)
		}

		fmt.Println(event.Object.Status.Phase)
	}
}

func startService(params Params, logger *zap.SugaredLogger) workflow.WorkflowService_WatchWorkflowsClient {
	url := params.ArgoServer
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		logger.Fatal(err)
	}

	client := workflow.NewWorkflowServiceClient(conn)
	request := &workflow.WatchWorkflowsRequest{}

	service, err := client.WatchWorkflows(context.TODO(), request)
	if err != nil {
		logger.Fatal(err)
	}
	return service
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
