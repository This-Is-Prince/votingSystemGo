package singleChoice

import (
	"testing"

	"github.com/thoas/go-funk"
)

type SingleChoiceVote struct {
	Choice  int       `json:"choice"`
	Balance float64   `json:"balance"`
	Scores  []float64 `json:"scores"`
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

func (v *SingleChoiceVoting) GetScoresTotal() float64 {
	return funk.Reduce(v.Votes, func(acc float64, vote SingleChoiceVote) float64 {
		return acc + vote.Balance
	}, float64(0)).(float64)
}

func (v *SingleChoiceVoting) GetScores(t *testing.T) []float64 {
	scores := []float64{}

	for range v.Choices {
		scores = append(scores, float64(0))
	}

	for _, vote := range v.Votes {
		choice := vote.Choice
		if IsValidChoice(choice, v.Choices) {
			scores[choice-1] = scores[choice-1] + vote.Balance
		}
	}

	return scores
}

func (v *SingleChoiceVoting) GetScoresByStrategy(t *testing.T) [][]float64 {
	scoresByStrategy := [][]float64{}

	for range v.Choices {
		scores := []float64{}
		for range v.Strategies {
			scores = append(scores, float64(0))
		}
		scoresByStrategy = append(scoresByStrategy, scores)
	}

	for _, vote := range v.Votes {
		choice := vote.Choice
		if IsValidChoice(choice, v.Choices) {
			for idx, score := range vote.Scores {
				scoresByStrategy[choice-1][idx] = scoresByStrategy[choice-1][idx] + score
			}
		}
	}

	return scoresByStrategy
}
