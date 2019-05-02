package genetic

type Population []*Species

func NewPopulation(pool *SpeciesPool, populationSize int) (result Population) {
	for i := 0; i < populationSize; i++ {
		result = append(result, pool.Get())
	}
	return
}

func (population Population) CloneFrom(leader *Species, freshBloodFraction float64) {
	freshBloodIdx := int(float64(len(population)) * (1 - freshBloodFraction))
	for _, species := range population[:freshBloodIdx] {
		if species.cloneOf != leader && species.cloneOf != leader.cloneOf {
			leader.CopyTo(species)
		}
	}
	for _, species := range population[freshBloodIdx:] {
		species.Reset()
	}
}
