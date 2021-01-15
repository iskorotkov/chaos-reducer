package reducer

import (
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/iskorotkov/chaos-reducer/api/metadata"
	"math/rand"
	"reflect"
	"testing"
)

const (
	minStages  = 1
	maxStages  = 10
	minSteps   = 1
	maxSteps   = 10
	iterations = 1000
)

func createStages(rng *rand.Rand) ([]Stage, int) {
	stages := createSucceededStages(rng)
	failedStageIndex := rng.Intn(len(stages))
	markStepsAsFailed(stages[failedStageIndex], rng)

	return stages, failedStageIndex
}

func markStepsAsFailed(stage Stage, rng *rand.Rand) {
	totalSteps := len(stage.Steps)
	indices := rng.Perm(totalSteps)

	var stepsFailed int
	if totalSteps == 1 {
		stepsFailed = 1
	} else {
		stepsFailed = 1 + rng.Intn(totalSteps-1)
	}

	for _, index := range indices[:stepsFailed] {
		stage.Steps[index].Phase = v1alpha1.NodeFailed
	}
}

func createSucceededStages(rng *rand.Rand) []Stage {
	stages := make([]Stage, 0)

	stagesNum := minStages + rng.Intn(maxStages-minStages)
	for i := 0; i < stagesNum; i++ {
		steps := make([]Step, 0)

		stepsNum := minSteps + rng.Intn(maxSteps-minSteps)
		for j := 0; j < stepsNum; j++ {
			steps = append(steps, NewStep(v1alpha1.NodeSucceeded, v1alpha1.WorkflowStep{}, v1alpha1.Template{}, metadata.TemplateMetadata{}))
		}

		stages = append(stages, NewStage(steps, []UtilityStep{}))
	}
	return stages
}

func Test_findFailedStage(t *testing.T) {
	rng := rand.New(rand.NewSource(0))

	for i := 0; i < iterations; i++ {
		stages, failedStageIndex := createStages(rng)
		stage, index := findFailedStage(stages)

		if failedStageIndex != index {
			t.Errorf("findFailedStage() = %v, want %v", index, failedStageIndex)
		}

		if !reflect.DeepEqual(stage, stages[failedStageIndex]) {
			t.Errorf("findFailedStage() = %v, want %v", stage, stages[failedStageIndex])
		}
	}
}

func Test_removeFollowingSteps(t *testing.T) {
	rng := rand.New(rand.NewSource(0))

	for i := 0; i < iterations; i++ {
		stages, failedStageIndex := createStages(rng)

		want := stages[:failedStageIndex]
		got := removeFollowingSteps(stages)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("removeFollowingSteps() = %v, want %v", got, want)
		}
	}
}

func Test_removePreviousSteps(t *testing.T) {
	rng := rand.New(rand.NewSource(0))

	for i := 0; i < iterations; i++ {
		stages, failedStageIndex := createStages(rng)

		want := stages[failedStageIndex:]
		got := removePreviousSteps(stages)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("removePreviousSteps() = %v, want %v", got, want)
		}
	}
}

func Test_removeRandomStages(t *testing.T) {
	rng := rand.New(rand.NewSource(0))

	containsStage := func(s Stage, ls []Stage) bool {
		for _, item := range ls {
			if reflect.DeepEqual(s, item) {
				return true
			}
		}

		return false
	}

	for i := 0; i < iterations; i++ {
		stages, failedStageIndex := createStages(rng)

		if len(stages) < 3 {
			i--
			continue
		}

		num := 2
		got := removeRandomStages(stages, num, rng)

		if len(got) > len(stages) {
			t.Errorf("removeRandomStages(), slice length = got %d, want <=%d", len(got), len(stages))
		}

		if len(got) < len(stages)-num {
			t.Errorf("removeRandomStages(), slice length = got %d, want >=%d", len(got), len(stages)-num)
		}

		if !containsStage(stages[failedStageIndex], got) {
			t.Errorf("removeRandomStages() does not contain failed stage")
		}

		for _, s := range got {
			if !containsStage(s, stages) {
				t.Errorf("removeRandomStages() contains stage not found in original list")
			}
		}
	}
}

func Test_removeRandomSteps(t *testing.T) {
	rng := rand.New(rand.NewSource(0))

	containsStep := func(s Step, ls []Step) bool {
		for _, item := range ls {
			if reflect.DeepEqual(s, item) {
				return true
			}
		}

		return false
	}

	failedSteps := func(steps []Step) int {
		cnt := 0
		for _, step := range steps {
			if step.Failed() {
				cnt++
			}
		}

		return cnt
	}

	totalFailedSteps := func(stages []Stage) int {
		cnt := 0
		for _, stage := range stages {
			cnt += failedSteps(stage.Steps)
		}

		return cnt
	}

	for i := 0; i < iterations; i++ {
		stages, _ := createStages(rng)

		num := 2
		got := removeRandomSteps(stages, num, rng)

		if len(got) != len(stages) {
			t.Errorf("removeRandomSteps(), slice length = got %d, want %d", len(got), len(stages))
		}

		if totalFailedSteps(got) != totalFailedSteps(stages) {
			t.Errorf("removeRandomSteps(), failed steps = got %d, want %d", totalFailedSteps(got), totalFailedSteps(stages))
		}

		for i, stage := range got {
			if len(stage.Steps) > len(stages[i].Steps) {
				t.Errorf("removeRandomSteps(), steps in stage = got %d, want <=%d", len(got), len(stages)-num)
			}

			if len(stage.Steps) < len(stages[i].Steps)-num {
				t.Errorf("removeRandomSteps(), steps in stage = got %d, want >=%d", len(got), len(stages)-num)
			}

			if len(stages[i].Steps) > num+failedSteps(stages[i].Steps) && len(stage.Steps) == len(stages[i].Steps) {
				t.Errorf("removeRandomSteps() mutated passed parameter instead of making a copy")
			}

			for _, step := range stage.Steps {
				if !containsStep(step, stages[i].Steps) {
					t.Errorf("removeRandomSteps() contains step not found in original list")
				}
			}
		}
	}
}
