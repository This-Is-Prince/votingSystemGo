package singleChoice

import (
	"math/big"
	"testing"
)

func TestSingleChoiceVoting(t *testing.T) {
	choices := []string{"First", "Second", "Third", "Fourth"}
	votes := []SingleChoiceVote{
		{
			Choice:  1,
			Balance: big.NewFloat(2.4946602468376033),
			Scores: []*big.Float{
				big.NewFloat(0.4946602468376035), big.NewFloat(2),
			},
		},
		{
			Choice:  2,
			Balance: big.NewFloat(0.4946602468376033),
			Scores: []*big.Float{
				big.NewFloat(2.4946602468376035), big.NewFloat(13),
			},
		},
		{
			Choice:  4,
			Balance: big.NewFloat(5.4946602468376033),
			Scores: []*big.Float{
				big.NewFloat(8.4946602468376035), big.NewFloat(22),
			},
		},
		{
			Choice:  3,
			Balance: big.NewFloat(2.2723898),
			Scores: []*big.Float{
				big.NewFloat(6.4946602468376035), big.NewFloat(5),
			},
		},
	}
	singleChoiceVoting := SingleChoiceVoting{
		Choices:    choices,
		Votes:      votes,
		Strategies: []interface{}{1, 2},
	}

	validVotes := singleChoiceVoting.GetValidVotes()
	if len(validVotes) != len(votes) {
		t.Errorf("Expected %d valid votes, got %d", len(votes), len(validVotes))
	}

	expectedScoresTotal := big.NewFloat(10.756370540512808).SetPrec(7)
	scoresTotal := singleChoiceVoting.GetScoresTotal().SetPrec(7)
	if scoresTotal.Cmp(expectedScoresTotal) != 0 {
		t.Errorf("Expected scores total to be %f, got %f", expectedScoresTotal, scoresTotal)
	}

	expectedScores := []*big.Float{
		big.NewFloat(2.4946602468376033),
		big.NewFloat(0.4946602468376033),
		big.NewFloat(2.2723898),
		big.NewFloat(5.494660246837603),
	}
	scores := singleChoiceVoting.GetScores(t)
	if len(scores) != len(choices) {
		t.Errorf("Expected %d scores, got %d", len(choices), len(scores))
	}

	for i, score := range scores {
		if score.SetPrec(5).Cmp(expectedScores[i].SetPrec(5)) != 0 {
			t.Errorf("Expected score %f for choice %s, got %f", expectedScores[i], choices[i], score)
		}
	}

	scoresByStrategy := singleChoiceVoting.GetScoresByStrategy(t)
	if len(scoresByStrategy) != len(choices) {
		t.Errorf("Expected %d scoresByStrategy, got %d", len(choices), len(scoresByStrategy))
	}

	expectedScoresByStrategy := [][]*big.Float{
		{big.NewFloat(0.4946602468376035), big.NewFloat(2)},
		{big.NewFloat(2.4946602468376033), big.NewFloat(13)},
		{big.NewFloat(6.494660246837603), big.NewFloat(5)},
		{big.NewFloat(8.494660246837604), big.NewFloat(22)},
	}

	for i, scoreByStrategy := range scoresByStrategy {
		for j, score := range scoreByStrategy {
			if score.SetPrec(5).Cmp(expectedScoresByStrategy[i][j].SetPrec(5)) != 0 {
				t.Errorf("Expected score %f got %f", expectedScoresByStrategy[i][j], score)
			}
		}
	}

}
