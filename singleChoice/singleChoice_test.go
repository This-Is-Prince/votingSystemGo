package singleChoice

import (
	"testing"

	"github.com/This-Is-Prince/votingSystemGo/utils"
)

func TestSingleChoiceVoting(t *testing.T) {
	choices := []string{"First", "Second", "Third", "Fourth"}
	votes := []SingleChoiceVote{
		{
			Choice:  1,
			Balance: float64(2.4946602468376033),
			Scores: []float64{
				float64(0.4946602468376035), float64(2),
			},
		},
		{
			Choice:  2,
			Balance: float64(0.4946602468376033),
			Scores: []float64{
				float64(2.4946602468376035), float64(13),
			},
		},
		{
			Choice:  4,
			Balance: float64(5.4946602468376033),
			Scores: []float64{
				float64(8.4946602468376035), float64(22),
			},
		},
		{
			Choice:  3,
			Balance: float64(2.2723898),
			Scores: []float64{
				float64(6.4946602468376035), float64(5),
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

	expectedScoresTotal := float64(10.756370540512808)
	scoresTotal := singleChoiceVoting.GetScoresTotal()
	if !utils.FloatEqual(scoresTotal, expectedScoresTotal) {
		t.Errorf("Expected scores total to be %f, got %f", expectedScoresTotal, scoresTotal)
	}

	expectedScores := []float64{
		float64(2.4946602468376033),
		float64(0.4946602468376033),
		float64(2.2723898),
		float64(5.494660246837603),
	}
	scores := singleChoiceVoting.GetScores(t)
	if len(scores) != len(choices) {
		t.Errorf("Expected %d scores, got %d", len(choices), len(scores))
	}

	for i, score := range scores {
		if !utils.FloatEqual(score, expectedScores[i]) {
			t.Errorf("Expected score %f for choice %s, got %f", expectedScores[i], choices[i], score)
		}
	}

	scoresByStrategy := singleChoiceVoting.GetScoresByStrategy(t)
	if len(scoresByStrategy) != len(choices) {
		t.Errorf("Expected %d scoresByStrategy, got %d", len(choices), len(scoresByStrategy))
	}

	expectedScoresByStrategy := [][]float64{
		{float64(0.4946602468376035), float64(2)},
		{float64(2.4946602468376033), float64(13)},
		{float64(6.494660246837603), float64(5)},
		{float64(8.494660246837604), float64(22)},
	}

	for i, scoreByStrategy := range scoresByStrategy {
		for j, score := range scoreByStrategy {
			if !utils.FloatEqual(score, expectedScoresByStrategy[i][j]) {
				t.Errorf("Expected score %f got %f", expectedScoresByStrategy[i][j], score)
			}
		}
	}

}
