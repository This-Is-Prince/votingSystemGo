package approval

import (
	"math/big"
	"testing"
)

func TestApprovalVoting(t *testing.T) {
	choices := []string{"First", "Second", "Third", "Fourth"}
	votes := []ApprovalVote{
		{
			Choice:  []int{4, 2, 3},
			Balance: big.NewFloat(2.4946602468376033),
			Scores:  []*big.Float{big.NewFloat(0.4946602468376035), big.NewFloat(2)},
		},
		{
			Choice:  []int{3, 1},
			Balance: big.NewFloat(12.812822710153798),
			Scores:  []*big.Float{big.NewFloat(10.812822710153798), big.NewFloat(2)},
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

	expectedScoresTotal := big.NewFloat(15.307483)
	scoresTotal := approvalVoting.GetScoresTotal()
	if scoresTotal.Cmp(expectedScoresTotal) != 0 {
		t.Errorf("Expected scores total to be %f, got %f", expectedScoresTotal, scoresTotal)
	}

	expectedScores := []*big.Float{big.NewFloat(12.812823), big.NewFloat(2.494660), big.NewFloat(15.307483), big.NewFloat(2.494660)}
	scores := approvalVoting.GetScores()
	if len(scores) != len(choices) {
		t.Errorf("Expected %d scores, got %d", len(choices), len(scores))
	}

	for i, score := range scores {
		if score.Cmp(expectedScores[i]) != 0 {
			t.Errorf("Expected score %f for choice %s, got %f", expectedScores[i], choices[i], score)
		}
	}

	scoresByStrategy := approvalVoting.GetScoresByStrategy()
	if len(scoresByStrategy) != len(choices) {
		t.Errorf("Expected %d scoresByStrategy, got %d", len(choices), len(scoresByStrategy))
	}

	expectedScoresByStrategy := [][]*big.Float{
		{big.NewFloat(10.812822710153798), big.NewFloat(2)},
		{big.NewFloat(0.4946602468376035), big.NewFloat(2)},
		{big.NewFloat(11.307482956991402), big.NewFloat(4)},
		{big.NewFloat(0.4946602468376035), big.NewFloat(2)},
	}

	for i, scoreByStrategy := range scoresByStrategy {
		for j, score := range scoreByStrategy {
			if score.Cmp(expectedScoresByStrategy[i][j]) != 0 {
				t.Errorf("Expected score %f got %f", expectedScoresByStrategy[i][j], score)
			}
		}
	}
}
