package weighted

import (
	"math"
	"testing"
)

func FloatEqual(a, b float64) bool {
	return math.Abs(a-b) < 0.0000001
}

func TestWeightedVoting(t *testing.T) {

	choices := []string{"First", "Second", "Third", "Fourth"}
	votes := []WeightedVote{
		{
			Choice: WeightedChoice{
				"1": 3,
				"2": 1,
				"3": 5,
				"4": 4,
			},
			Balance: 2.4946602468376033,
			Scores:  []float64{(0.4946602468376035), 2.0},
		},
		{
			Choice: WeightedChoice{
				"1": 2,
				"2": 3,
				"3": 6,
				"4": 1,
			},
			Balance: (0.4946602468376033),
			Scores: []float64{
				(2.4946602468376035), (13),
			},
		},
		{
			Choice: WeightedChoice{
				"1": 5,
				"2": 5,
				"3": 1,
				"4": 8,
			},
			Balance: (5.4946602468376033),
			Scores: []float64{
				(8.4946602468376035), (22),
			},
		},
		{
			Choice: WeightedChoice{
				"1": 9,
				"2": 2,
				"3": 4,
				"4": 5,
			},
			Balance: (2.2723898),
			Scores: []float64{
				(6.4946602468376035), (5),
			},
		},
	}
	quadraticVoting := WeightedVoting{
		Choices:    choices,
		Votes:      votes,
		Strategies: []interface{}{1, 2},
	}

	validVotes := quadraticVoting.GetValidVotes()
	if len(validVotes) != len(votes) {
		t.Errorf("Expected %d valid votes, got %d", len(votes), len(validVotes))
	}

	expectedScoresTotal := (10.756370540512808)
	scoresTotal := quadraticVoting.GetScoresTotal()
	if !FloatEqual(scoresTotal, expectedScoresTotal) {
		t.Errorf("Expected scores total to be %f, got %f", expectedScoresTotal, scoresTotal)
	}

	expectedScores := []float64{
		(3.1266728335182266),
		(1.9887642066258324),
		(1.9504854383113572),
		(3.690448062057391),
	}
	scores := quadraticVoting.GetScores()
	if len(scores) != len(choices) {
		t.Errorf("Expected %d scores, got %d", len(choices), len(scores))
	}

	for i, score := range scores {
		if !FloatEqual(score, expectedScores[i]) {
			t.Errorf("Expected score %f for choice %s, got %f", expectedScores[i], choices[i], score)
		}
	}

	scoresByStrategy := quadraticVoting.GetScoresByStrategy()
	if len(scoresByStrategy) != len(choices) {
		t.Errorf("Expected %d scoresByStrategy, got %d", len(choices), len(scoresByStrategy))
	}

	expectedScoresByStrategy := [][]float64{
		{(1.0200604351166154), (1.9131061361947228)},
		{(0.6360388490574986), (1.7383678294383311)},
		{(0.5709368916283944), (1.6906294208404111)},
		{(0.9971936682110845), (2.1900373100257506)},
	}

	for i, scoreByStrategy := range scoresByStrategy {
		for j, score := range scoreByStrategy {
			if !FloatEqual(score, expectedScoresByStrategy[i][j]) {
				t.Errorf("Expected score %f got %f for %v %v", expectedScoresByStrategy[i][j], score, i, j)
			}
		}
	}
}
