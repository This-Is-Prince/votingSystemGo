package quadratic

import (
	"math/big"
	"testing"
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
			Balance: big.NewFloat(2.4946602468376033),
			Scores:  []*big.Float{big.NewFloat(0.4946602468376035), big.NewFloat(2)},
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

	expectedScoresTotal := big.NewFloat(2.494660).SetPrec(7)
	scoresTotal := quadraticVoting.GetScoresTotal().SetPrec(7)
	if scoresTotal.Cmp(expectedScoresTotal) != 0 {
		t.Errorf("Expected scores total to be %f, got %f", expectedScoresTotal, scoresTotal)
	}

	expectedScores := []*big.Float{
		big.NewFloat(0.7483980740512811),
		big.NewFloat(0.2494660246837603),
		big.NewFloat(0.9978640987350412),
		big.NewFloat(0.4989320493675206),
	}
	scores := quadraticVoting.GetScores()
	if len(scores) != len(choices) {
		t.Errorf("Expected %d scores, got %d", len(choices), len(scores))
	}

	for i, score := range scores {
		if score.SetPrec(5).Cmp(expectedScores[i].SetPrec(5)) != 0 {
			t.Errorf("Expected score %f for choice %s, got %f", expectedScores[i], choices[i], score)
		}
	}

	scoresByStrategy := quadraticVoting.GetScoresByStrategy()
	if len(scoresByStrategy) != len(choices) {
		t.Errorf("Expected %d scoresByStrategy, got %d", len(choices), len(scoresByStrategy))
	}

	expectedScoresByStrategy := [][]*big.Float{
		{big.NewFloat(0.14808111041719751), big.NewFloat(0.598718459241025)},
		{big.NewFloat(0.049360370139065836), big.NewFloat(0.19957281974700825)},
		{big.NewFloat(0.19744148055626334), big.NewFloat(0.798291278988033)},
		{big.NewFloat(0.09872074027813167), big.NewFloat(0.39914563949401655)},
	}

	for i, scoreByStrategy := range scoresByStrategy {
		for j, score := range scoreByStrategy {
			if score.SetPrec(5).Cmp(expectedScoresByStrategy[i][j].SetPrec(5)) != 0 {
				t.Errorf("Expected score %f got %f", expectedScoresByStrategy[i][j], score)
			}
		}
	}

}
