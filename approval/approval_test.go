package approval

import (
	"testing"

	"github.com/This-Is-Prince/votingSystemGo/utils"
)

func TestApprovalVoting(t *testing.T) {
	choices := []string{"First", "Second", "Third", "Fourth"}
	votes := []ApprovalVote{
		{
			Choice:  []int{4, 2, 3},
			Balance: float64(2.4946602468376033),
			Scores:  []float64{float64(0.4946602468376035), float64(2)},
		},
		{
			Choice:  []int{3, 1},
			Balance: float64(12.812822710153798),
			Scores:  []float64{float64(10.812822710153798), float64(2)},
		},
	}
	strategies := []interface{}{1, 2}
	approvalVoting := ApprovalVoting{
		Choices:    choices,
		Votes:      votes,
		Strategies: strategies,
	}

	validVotes := approvalVoting.GetValidVotes()
	if len(validVotes) != len(votes) {
		t.Errorf("Expected %d valid votes, got %d", len(votes), len(validVotes))
	}

	expectedScoresTotal := float64(15.307483)
	scoresTotal := approvalVoting.GetScoresTotal()
	if !utils.FloatEqual(scoresTotal, expectedScoresTotal) {
		t.Errorf("Expected scores total to be %f, got %f", expectedScoresTotal, scoresTotal)
	}

	expectedScores := []float64{float64(12.812822710153798), float64(2.4946602468376033), float64(15.307483), float64(2.4946602468376033)}
	scores := approvalVoting.GetScores()
	if len(scores) != len(choices) {
		t.Errorf("Expected %d scores, got %d", len(choices), len(scores))
	}

	for i, score := range scores {
		if !utils.FloatEqual(score, expectedScores[i]) {
			t.Errorf("Expected score %f for choice %s, got %f", expectedScores[i], choices[i], score)
		}
	}

	scoresByStrategy := approvalVoting.GetScoresByStrategy()
	if len(scoresByStrategy) != len(choices) {
		t.Errorf("Expected %d scoresByStrategy, got %d", len(choices), len(scoresByStrategy))
	}

	expectedScoresByStrategy := [][]float64{
		{float64(10.812822710153798), float64(2)},
		{float64(0.4946602468376035), float64(2)},
		{float64(11.307482956991402), float64(4)},
		{float64(0.4946602468376035), float64(2)},
	}

	for i, scoreByStrategy := range scoresByStrategy {
		for j, score := range scoreByStrategy {
			if !utils.FloatEqual(score, expectedScoresByStrategy[i][j]) {
				t.Errorf("Expected score %f got %f", expectedScoresByStrategy[i][j], score)
			}
		}
	}
}
