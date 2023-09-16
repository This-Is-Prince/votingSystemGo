package singleChoice

import (
	"math/big"
	"testing"

	"github.com/thoas/go-funk"
)

type SingleChoiceVote struct {
	Choice  int          `json:"choice"`
	Balance *big.Float   `json:"balance"`
	Scores  []*big.Float `json:"scores"`
}

type SingleChoiceVoting struct {
	Choices    []string           `json:"choices"`
	Votes      []SingleChoiceVote `json:"votes"`
	Strategies []interface{}      `json:"strategies"`
}

func IsValidChoice(voteChoice int, proposalChoices []string) bool {
	return voteChoice > 0 && voteChoice <= len(proposalChoices)
}

func (v *SingleChoiceVoting) GetValidVotes() []SingleChoiceVote {
	return funk.Filter(v.Votes, func(vote SingleChoiceVote) bool {
		return IsValidChoice(vote.Choice, v.Choices)
	}).([]SingleChoiceVote)
}

func (v *SingleChoiceVoting) GetScoresTotal() *big.Float {
	return funk.Reduce(v.Votes, func(acc *big.Float, vote SingleChoiceVote) *big.Float {
		return acc.Add(acc, vote.Balance)
	}, big.NewFloat(0)).(*big.Float)
}

func (v *SingleChoiceVoting) GetScores(t *testing.T) []*big.Float {
	scores := []*big.Float{}

	for range v.Choices {
		scores = append(scores, big.NewFloat(0))
	}

	for _, vote := range v.Votes {
		choice := vote.Choice
		if IsValidChoice(choice, v.Choices) {
			scores[choice-1] = scores[choice-1].Add(scores[choice-1], vote.Balance)
		}
	}

	return scores
}

func (v *SingleChoiceVoting) GetScoresByStrategy(t *testing.T) [][]*big.Float {
	scoresByStrategy := [][]*big.Float{}

	for range v.Choices {
		scores := []*big.Float{}
		for range v.Strategies {
			scores = append(scores, big.NewFloat(0))
		}
		scoresByStrategy = append(scoresByStrategy, scores)
	}

	for _, vote := range v.Votes {
		choice := vote.Choice
		if IsValidChoice(choice, v.Choices) {
			for idx, score := range vote.Scores {
				scoresByStrategy[choice-1][idx] = scoresByStrategy[choice-1][idx].Add(scoresByStrategy[choice-1][idx], score)
			}
		}
	}

	return scoresByStrategy
}
