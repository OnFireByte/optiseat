package optimize

import (
	"math"
	"math/rand"
	"slices"
)

const (
	initialTemp   = 200.
	coolingRate   = 0.995
	maxIterations = 50000
)

type Pair struct {
	A int
	B int
}

// Helper function to calculate the happiness score
func calculateHappiness(tables [][]int, pref [][]int) int {
	happiness := 0
	for _, table := range tables {
		n := len(table)
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				// Pairs on the same table
				happiness += pref[table[i]][table[j]]
				happiness += pref[table[j]][table[i]]
			}
			// Adjacent pairs on the round table
			happiness += pref[table[i]][table[(i+1)%n]]
			happiness += pref[table[(i+1)%n]][table[i]]
		}
	}
	return happiness
}

// Helper function to generate an initial random seating arrangement
func initialSeating(n, m int) [][]int {
	people := rand.Perm(n)
	tables := [][]int{}
	for i := 0; i < n; i += m {
		end := i + m
		if end > n {
			end = n
		}
		tables = append(tables, people[i:end])
	}
	return tables
}

// Helper function to generate a neighboring seating arrangement
func neighbor(tables [][]int, max_seat int) [][]int {
	newTables := make([][]int, len(tables))
	for i := range tables {
		newTables[i] = append([]int(nil), tables[i]...)
	}

	allSeated := true
	for _, table := range tables {
		allSeated = allSeated && len(table) == max_seat
	}

	if rand.Float64() < 0.5 || allSeated {
		// Swap two people between tables
		t1, t2 := rand.Intn(len(newTables)), rand.Intn(len(newTables))
		for len(newTables[t1]) <= 0 || len(newTables[t2]) <= 0 {
			t1, t2 = rand.Intn(len(newTables)), rand.Intn(len(newTables))
		}
		i := rand.Intn(len(newTables[t1]))
		j := rand.Intn(len(newTables[t2]))
		newTables[t1][i], newTables[t2][j] = newTables[t2][j], newTables[t1][i]
	} else {
		// Move one person to a different table
		t1, t2 := rand.Intn(len(newTables)), rand.Intn(len(newTables))
		for len(newTables[t1]) <= 0 || len(newTables[t2]) >= max_seat || t1 == t2 {
			t1, t2 = rand.Intn(len(newTables)), rand.Intn(len(newTables))
		}
		i := rand.Intn(len(newTables[t1]))

		person := newTables[t1][i]
		newTables[t1] = slices.Delete(newTables[t1], i, i+1)
		newTables[t2] = append(newTables[t2], person)
	}

	return newTables
}

// Simulated Annealing Algorithm.
//
// pref is a 2D slice that represent a person's preference A to sit with person B
// in the way that pref[A][B] = score.
//
// n is the number of peoples
//
// maxSeat is number of maximum amount of people that can sit in the table
func SimulatedAnnealing(pref [][]int, n int, maxSeat int) ([][]int, int) {
	currentSeating := initialSeating(n, maxSeat)
	currentScore := calculateHappiness(currentSeating, pref)
	bestSeating := currentSeating
	bestScore := currentScore
	temperature := initialTemp
	for iteration := 0; iteration < maxIterations; iteration++ {
		newSeating := neighbor(currentSeating, maxSeat)
		newScore := calculateHappiness(newSeating, pref)

		// Accept the new seating with probability based on temperature
		if newScore > currentScore || rand.Float64() < math.Exp(float64(newScore-currentScore)/temperature) {
			// if newScore > currentScore {
			currentSeating = newSeating
			currentScore = newScore
		}

		// Update the best found solution
		if currentScore > bestScore {
			bestSeating = currentSeating
			bestScore = currentScore
		}

		// Cool down the temperature
		temperature *= coolingRate
	}

	return bestSeating, bestScore
}
