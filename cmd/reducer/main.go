package main

import (
	"github.com/caarlos0/env"
	"go.uber.org/zap"
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

	listenToEvents(params, sugar)
}

func listenToEvents(params Params, logger *zap.SugaredLogger) {
	// service, err := client.WatchWorkflows(context.Background(), &workflow.WatchWorkflowsRequest{})
	// if err != nil {
	// 	logger.Fatal(err)
	// }
	//
	// for {
	// 	event, err := service.Recv()
	// 	if err != nil {
	// 		logger.Fatal(err)
	// 	}
	//
	// 	wf := event.Object
	// 	status := wf.Status
	//
	// 	// Skip workflows that hasn't failed
	// 	if status.Phase != "Failed" {
	// 		continue
	// 	}
	//
	// 	// Skip workflows that finished more than 60 seconds ago
	// 	// if !status.FinishedAt.Add(time.Second * 60).After(time.Now()) {
	// 	// 	continue
	// 	// }
	//
	// 	fmt.Printf("%s: %s (%s - %s)\n", wf.Name, status.Phase, status.StartedAt.Format("02.01.2006 15:04:05 MST"), status.FinishedAt.Format("02.01.2006 15:04:05 MST"))
	//
	// 	_, ok := wf.Annotations["chaosframework.com/next"]
	// 	if ok {
	// 		// Workflow already has follow-up workflow
	// 		continue
	// 	}
	//
	// 	iteration, ok := wf.Annotations["chaosframework.com/iteration"]
	// 	if !ok {
	// 		// Workflow wasn't launched as a part of test
	// 		continue
	// 	}
	//
	// 	iterationNum, err := strconv.ParseInt(iteration, 10, 32)
	// 	if err != nil {
	// 		logger.Fatal(err)
	// 	}
	//
	// 	wf.Annotations["chaosframework.com/iteration"] = strconv.FormatInt(iterationNum+1, 10)
	//
	// 	// TODO: launch new workflow
	//
	// 	// TODO: assign new workflow name
	// 	name := "name"
	// 	namespace := "namespace"
	//
	// 	wf.Annotations["chaosframework.com/next"] = fmt.Sprintf("%s/%s", namespace, name)
	// }
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
