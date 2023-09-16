package approval

import (
	"math/big"

	"github.com/thoas/go-funk"
)

type ApprovalVote struct {
	Choice  []int        `json:"choice"`
	Balance *big.Float   `json:"balance"`
	Scores  []*big.Float `json:"scores"`
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

func (v *ApprovalVoting) GetScoresTotal() *big.Float {
	return funk.Reduce(v.Votes, func(acc *big.Float, vote ApprovalVote) *big.Float {
		return acc.Add(acc, vote.Balance)
	}, big.NewFloat(0)).(*big.Float)
}

func (v *ApprovalVoting) GetScores() []*big.Float {
	scores := []*big.Float{}

	for range v.Choices {
		scores = append(scores, big.NewFloat(0))
	}

	for _, vote := range v.Votes {
		if IsValidChoice(vote.Choice, v.Choices) {
			for _, choice := range vote.Choice {
				scores[choice-1] = scores[choice-1].Add(scores[choice-1], vote.Balance)
			}
		}
	}

	return scores
}

func (v *ApprovalVoting) GetScoresByStrategy() [][]*big.Float {
	scoresByStrategy := [][]*big.Float{}

	for range v.Choices {
		scores := []*big.Float{}
		for range v.Strategies {
			scores = append(scores, big.NewFloat(0))
		}
		scoresByStrategy = append(scoresByStrategy, scores)
	}

	for _, vote := range v.Votes {
		if IsValidChoice(vote.Choice, v.Choices) {
			for _, choice := range vote.Choice {
				for idx, score := range vote.Scores {
					scoresByStrategy[choice-1][idx] = scoresByStrategy[choice-1][idx].Add(scoresByStrategy[choice-1][idx], score)
				}
			}
		}
	}

	return scoresByStrategy
}
