package reducer

import (
	"math/rand"
)

func Reduce(scenario Scenario) (Scenario, error) {
	panic("not implemented")
}

func removeFollowingSteps(stages []Stage) []Stage {
	_, index := findFailedStage(stages)
	return stages[0:index]
}

func stageFailed(stage Stage) bool {
	for _, step := range stage.Steps {
		if step.Failed {
			return true
		}
	}

	return false
}

func findFailedStage(stages []Stage) (Stage, int) {
	for index, stage := range stages {
		if stageFailed(stage) {
			return stage, index
		}
	}

	panic("no failed stage found")
}

func removePreviousSteps(stages []Stage) []Stage {
	_, index := findFailedStage(stages)
	return stages[index:]
}

func removeRandomStages(stages []Stage, num int, rng *rand.Rand) []Stage {
	// Make a copy of stages slice
	ls := make([]Stage, len(stages))
	copy(ls, stages)
	stages = ls

	if num > len(stages)-1 {
		panic("can't remove more stages than present")
	}

	for i := 0; i < num; i++ {
		index := rng.Intn(len(stages))
		for stageFailed(stages[index]) {
			index = rng.Intn(len(stages))
		}

		stages = append(stages[:index], stages[index+1:]...)
	}

	return stages
}

func removeRandomSteps(stages []Stage, maxStepsPerStage int, rng *rand.Rand) []Stage {
	// Copy all stages in slice
	ls := make([]Stage, len(stages))
	for i, s := range stages {
		ls[i] = s

		// Copy steps slice
		ls[i].Steps = make([]Step, len(s.Steps))
		copy(ls[i].Steps, s.Steps)
	}
	stages = ls

	for i, stage := range stages {
		for j := 0; j < maxStepsPerStage; j++ {
			if len(stage.Steps) == 1 {
				continue
			}

			if len(stage.Steps) == failedSteps(stage.Steps) {
				continue
			}

			index := rng.Intn(len(stage.Steps))
			for stage.Steps[index].Failed {
				index = rng.Intn(len(stage.Steps))
			}

			stage.Steps = append(stage.Steps[:index], stage.Steps[index+1:]...)
		}

		stages[i] = stage
	}

	return stages
}

func failedSteps(steps []Step) int {
	cnt := 0
	for _, step := range steps {
		if step.Failed {
			cnt++
		}
	}

	return cnt
}
