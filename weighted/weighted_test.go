package weighted

import (
	"math/big"
	"testing"
)

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
			Balance: big.NewFloat(2.4946602468376033),
			Scores:  []*big.Float{big.NewFloat(0.4946602468376035), big.NewFloat(2)},
		},
		{
			Choice: WeightedChoice{
				"1": 2,
				"2": 3,
				"3": 6,
				"4": 1,
			},
			Balance: big.NewFloat(0.4946602468376033),
			Scores: []*big.Float{
				big.NewFloat(2.4946602468376035), big.NewFloat(13),
			},
		},
		{
			Choice: WeightedChoice{
				"1": 5,
				"2": 5,
				"3": 1,
				"4": 8,
			},
			Balance: big.NewFloat(5.4946602468376033),
			Scores: []*big.Float{
				big.NewFloat(8.4946602468376035), big.NewFloat(22),
			},
		},
		{
			Choice: WeightedChoice{
				"1": 9,
				"2": 2,
				"3": 4,
				"4": 5,
			},
			Balance: big.NewFloat(2.2723898),
			Scores: []*big.Float{
				big.NewFloat(6.4946602468376035), big.NewFloat(5),
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

	expectedScoresTotal := big.NewFloat(10.756370540512808).SetPrec(7)
	scoresTotal := quadraticVoting.GetScoresTotal().SetPrec(7)
	if scoresTotal.Cmp(expectedScoresTotal) != 0 {
		t.Errorf("Expected scores total to be %f, got %f", expectedScoresTotal, scoresTotal)
	}

	expectedScores := []*big.Float{
		big.NewFloat(3.1266728335182266),
		big.NewFloat(1.9887642066258324),
		big.NewFloat(1.9504854383113572),
		big.NewFloat(3.690448062057391),
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
		{big.NewFloat(1.0200604351166154), big.NewFloat(1.9131061361947228)},
		{big.NewFloat(0.6360388490574986), big.NewFloat(1.7383678294383311)},
		{big.NewFloat(0.5709368916283944), big.NewFloat(1.6906294208404111)},
		{big.NewFloat(0.9971936682110845), big.NewFloat(2.1900373100257506)},
	}

	for i, scoreByStrategy := range scoresByStrategy {
		for j, score := range scoreByStrategy {
			if score.SetPrec(5).Cmp(expectedScoresByStrategy[i][j].SetPrec(5)) != 0 {
				t.Errorf("Expected score %f got %f for %v %v", expectedScoresByStrategy[i][j], score, i, j)
			}
		}
	}
}
