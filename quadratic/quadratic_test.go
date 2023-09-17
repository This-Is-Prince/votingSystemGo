package quadratic

import (
	"testing"

	"github.com/This-Is-Prince/votingSystemGo/utils"
)

func TestQuadraticVoting(t *testing.T) {
	choices := []string{"First", "Second", "Third", "Fourth"}
	votes := []QuadraticVote{
		{
			Choice: QuadraticChoice{
				"1": 3,
				"2": 1,
				"3": 4,
				"4": 2,
			},
			Balance: float64(2.4946602468376033),
			Scores:  []float64{float64(0.4946602468376035), float64(2)},
		},
	}
	quadraticVoting := QuadraticVoting{
		Choices:    choices,
		Votes:      votes,
		Strategies: []interface{}{1, 2},
	}

	validVotes := quadraticVoting.GetValidVotes()
	if len(validVotes) != len(votes) {
		t.Errorf("Expected %d valid votes, got %d", len(votes), len(validVotes))
	}

	expectedScoresTotal := float64(2.494660)
	scoresTotal := quadraticVoting.GetScoresTotal()
	if !utils.FloatEqual(scoresTotal, expectedScoresTotal) {
		t.Errorf("Expected scores total to be %f, got %f", expectedScoresTotal, scoresTotal)
	}

	expectedScores := []float64{
		float64(0.7483980740512811),
		float64(0.2494660246837603),
		float64(0.9978640987350412),
		float64(0.4989320493675206),
	}
	scores := quadraticVoting.GetScores()
	if len(scores) != len(choices) {
		t.Errorf("Expected %d scores, got %d", len(choices), len(scores))
	}

	for i, score := range scores {
		if !utils.FloatEqual(score, expectedScores[i]) {
			t.Errorf("Expected score %f for choice %s, got %f", expectedScores[i], choices[i], score)
		}
	}

	scoresByStrategy := quadraticVoting.GetScoresByStrategy()
	if len(scoresByStrategy) != len(choices) {
		t.Errorf("Expected %d scoresByStrategy, got %d", len(choices), len(scoresByStrategy))
	}

	expectedScoresByStrategy := [][]float64{
		{float64(0.14808111041719751), float64(0.598718459241025)},
		{float64(0.049360370139065836), float64(0.19957281974700825)},
		{float64(0.19744148055626334), float64(0.798291278988033)},
		{float64(0.09872074027813167), float64(0.39914563949401655)},
	}

	for i, scoreByStrategy := range scoresByStrategy {
		for j, score := range scoreByStrategy {
			if !utils.FloatEqual(score, expectedScoresByStrategy[i][j]) {
				t.Errorf("Expected score %f got %f", expectedScoresByStrategy[i][j], score)
			}
		}
	}

}
