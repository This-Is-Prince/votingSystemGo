package approval

import (
	"github.com/thoas/go-funk"
)

type ApprovalVote struct {
	Choice  []int     `json:"choice"`
	Balance float64   `json:"balance"`
	Scores  []float64 `json:"scores"`
}

type ApprovalVoting struct {
	Choices    []string       `json:"choices"`
	Votes      []ApprovalVote `json:"votes"`
	Strategies []interface{}  `json:"strategies"`
}

func IsValidChoice(voteChoice []int, proposalChoices []string) bool {
	voteChoiceSet := make(map[int]struct{})
	filteredVoteChoice := funk.FilterInt(voteChoice, func(c int) bool {
		voteChoiceSet[c] = struct{}{}
		return c > 0 && c <= len(proposalChoices)
	})

	return len(voteChoice) == len(filteredVoteChoice) && len(voteChoice) == len(voteChoiceSet)
}

func (v *ApprovalVoting) GetValidVotes() []ApprovalVote {
	return funk.Filter(v.Votes, func(vote ApprovalVote) bool {
		return IsValidChoice(vote.Choice, v.Choices)
	}).([]ApprovalVote)
}

func (v *ApprovalVoting) GetScoresTotal() float64 {
	return funk.Reduce(v.Votes, func(acc float64, vote ApprovalVote) float64 {
		return acc + vote.Balance
	}, float64(0)).(float64)
}

func (v *ApprovalVoting) GetScores() []float64 {
	scores := []float64{}

	for range v.Choices {
		scores = append(scores, float64(0))
	}

	for _, vote := range v.Votes {
		if IsValidChoice(vote.Choice, v.Choices) {
			for _, choice := range vote.Choice {
				scores[choice-1] = scores[choice-1] + vote.Balance
			}
		}
	}

	return scores
}

func (v *ApprovalVoting) GetScoresByStrategy() [][]float64 {
	scoresByStrategy := [][]float64{}

	for range v.Choices {
		scores := []float64{}
		for range v.Strategies {
			scores = append(scores, float64(0))
		}
		scoresByStrategy = append(scoresByStrategy, scores)
	}

	for _, vote := range v.Votes {
		if IsValidChoice(vote.Choice, v.Choices) {
			for _, choice := range vote.Choice {
				for idx, score := range vote.Scores {
					scoresByStrategy[choice-1][idx] = scoresByStrategy[choice-1][idx] + score
				}
			}
		}
	}

	return scoresByStrategy
}
