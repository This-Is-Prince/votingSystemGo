package weighted

import (
	"log"
	"strconv"
	"testing"

	"github.com/thoas/go-funk"
)

type WeightedVote struct {
	Choice  WeightedChoice `json:"choice"`
	Balance float64        `json:"balance"`
	Scores  []float64      `json:"scores"`
}

type WeightedChoice map[string]int

type WeightedVoting struct {
	Choices    []string       `json:"choices"`
	Votes      []WeightedVote `json:"votes"`
	Strategies []interface{}  `json:"strategies"`
}

func IsValidChoice(voteChoice WeightedChoice, proposalChoices []string) bool {
	if voteChoice == nil || len(voteChoice) == 0 {
		return false
	}

	for k, v := range voteChoice {
		if v < 0 {
			return false
		}

		numKey, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			return false
		}

		if numKey <= 0 || int(numKey) > len(proposalChoices) {
			return false
		}
	}

	return true
}

func CalcPercentageOfSum(choice float64, choices []float64) float64 {
	if choice == 0.0 {
		return 0.0
	}

	whole := funk.Reduce(choices, func(acc float64, c float64) float64 {
		return acc + c
	}, 0).(float64)

	if whole == 0.0 {
		return 0.0
	}

	return choice / whole
}

func WeightedPower(choice float64, choices []float64, balance float64) float64 {
	percentage := CalcPercentageOfSum(choice, choices)
	return percentage * balance
}

func CalcReducedQuadraticScores(scoresTotal float64, percentages []float64) []float64 {
	return funk.Map(percentages, func(p float64) float64 {
		return p * scoresTotal
	}).([]float64)
}

func (v *WeightedVoting) GetValidVotes() []WeightedVote {
	return funk.Filter(v.Votes, func(vote WeightedVote) bool {
		return IsValidChoice(vote.Choice, v.Choices)
	}).([]WeightedVote)
}

func (v *WeightedVoting) GetScoresTotal() float64 {
	return funk.Reduce(v.Votes, func(acc float64, vote WeightedVote) float64 {
		return acc + vote.Balance
	}, 0).(float64)
}

func (v *WeightedVoting) GetScores(t *testing.T) []float64 {
	scoresTotal := 0.0
	scores := []float64{}

	for range v.Choices {
		scores = append(scores, 0.0)
	}

	for _, vote := range v.Votes {
		if IsValidChoice(vote.Choice, v.Choices) {
			scoresTotal = scoresTotal + vote.Balance
			choices := []float64{}
			for _, v := range vote.Choice {
				choices = append(choices, (float64(v)))
			}

			for idx, value := range vote.Choice {
				choiceWeightedPower := WeightedPower((float64(value)), choices, vote.Balance)
				index, err := strconv.ParseInt(idx, 10, 64)
				if err != nil {
					log.Println("Error while parsing string:-", err)
					continue
				}
				scores[index-1] = scores[index-1] + choiceWeightedPower
			}

		}
	}

	percentageOfScores := []float64{}
	for _, score := range scores {
		percentageOfScores = append(percentageOfScores, CalcPercentageOfSum(score, scores))
	}
	newScores := CalcReducedQuadraticScores(scoresTotal, percentageOfScores)

	t.Error("NewScores:-", newScores)

	return newScores
}

func (v *WeightedVoting) GetScoresByStrategy(t *testing.T) [][]float64 {
	scoresTotal := 0.0
	scoresByStrategy := [][]float64{}

	for range v.Choices {
		scores := []float64{}
		for range v.Strategies {
			scores = append(scores, 0.0)
		}
		scoresByStrategy = append(scoresByStrategy, scores)
	}

	for _, vote := range v.Votes {
		if IsValidChoice(vote.Choice, v.Choices) {
			scoresTotal = (scoresTotal + vote.Balance)
			choices := []float64{}
			for _, v := range vote.Choice {
				choices = append(choices, (float64(v)))
			}
			for idx, value := range vote.Choice {
				index, err := strconv.ParseInt(idx, 10, 64)
				if err != nil {
					log.Println("Error while parsing string:-", err)
					continue
				}
				for sIdx, score := range vote.Scores {
					choiceWeightedPower := WeightedPower((float64(value)), choices, score)
					scoresByStrategy[index-1][sIdx] = scoresByStrategy[index-1][sIdx] + choiceWeightedPower
				}
			}
		}
	}

	flattenScoresByStrategy := funk.FlattenDeep(scoresByStrategy)

	for idx, scores := range scoresByStrategy {
		percentageOfScores := []float64{}
		for _, score := range scores {
			percentageOfScores = append(percentageOfScores, CalcPercentageOfSum((score), flattenScoresByStrategy.([]float64)))
		}
		scoresByStrategy[idx] = CalcReducedQuadraticScores(scoresTotal, percentageOfScores)
	}

	return scoresByStrategy
}
