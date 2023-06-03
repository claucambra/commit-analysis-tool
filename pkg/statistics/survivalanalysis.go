package statistics

type TimeStepPopulation []int
type TimeStepSurvival []float64

func (tsp *TimeStepPopulation) KaplanMeierSurvival() TimeStepSurvival {
	tspLen := len(*tsp)
	survival := make(TimeStepSurvival, tspLen)

	if tspLen == 0 {
		return survival
	}

	prevTimeStepSurvival := 1.0
	prevTimeStepPopulation := (*tsp)[0]

	for timestep, population := range *tsp {
		events := prevTimeStepPopulation - population
		timeStepSurvival := prevTimeStepSurvival * (1 - (float64(events) / float64(prevTimeStepPopulation)))

		survival[timestep] = timeStepSurvival

		prevTimeStepSurvival = timeStepSurvival
		prevTimeStepPopulation = population
	}

	return survival
}
