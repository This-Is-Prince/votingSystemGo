package quadratic

import (
	"log"
	"math"
	"strconv"

	"github.com/thoas/go-funk"

	"github.com/This-Is-Prince/votingSystemGo/utils"
)

type QuadraticVote struct {
	Choice  QuadraticChoice `json:"choice"`
	Balance float64         `json:"balance"`
	Scores  []float64       `json:"scores"`
}

type QuadraticChoice map[string]int

type QuadraticVoting struct {
	Choices    []string        `json:"choices"`
	Votes      []QuadraticVote `json:"votes"`
	Strategies []interface{}   `json:"strategies"`
}

func IsValidChoice(voteChoice QuadraticChoice, proposalChoices []string) bool {
	if voteChoice == nil || len(voteChoice) == 0 {
		return false
	}

	for k, v := range voteChoice {
		if v <= 0 || v > len(proposalChoices) {
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

func (v *QuadraticVoting) GetValidVotes() []QuadraticVote {
	return funk.Filter(v.Votes, func(vote QuadraticVote) bool {
		return IsValidChoice(vote.Choice, v.Choices)
	}).([]QuadraticVote)
}

func (v *QuadraticVoting) GetScoresTotal() float64 {
	return funk.Reduce(v.Votes, func(acc float64, vote QuadraticVote) float64 {
		return acc + vote.Balance
	}, float64(0)).(float64)
}

func (v *QuadraticVoting) GetScores() []float64 {
	scoresTotal := float64(0)
	scores := []float64{}

	for range v.Choices {
		scores = append(scores, float64(0))
	}

	for _, vote := range v.Votes {
		if IsValidChoice(vote.Choice, v.Choices) {
			scoresTotal = scoresTotal + vote.Balance
			choices := []float64{}
			for _, v := range vote.Choice {
				choices = append(choices, float64(v))
			}
			for idx, value := range vote.Choice {
				choiceWeightPercent := utils.CalcPercentageOfSum(float64(value), choices)
				choiceWeightPower := choiceWeightPercent * vote.Balance
				sqrt := math.Sqrt(choiceWeightPower)
				index, err := strconv.ParseInt(idx, 10, 64)
				if err != nil {
					log.Println("Error while parsing string:-", err)
					continue
				}
				scores[index-1] = scores[index-1] + sqrt
			}
		}
	}

	for idx, score := range scores {
		scores[idx] = score * score
	}

	percentageOfScores := []float64{}
	for _, score := range scores {
		percentageOfScores = append(percentageOfScores, utils.CalcPercentageOfSum(score, scores))
	}
	return utils.CalcReducedQuadraticScores(scoresTotal, percentageOfScores)
}

func (v *QuadraticVoting) GetScoresByStrategy() [][]float64 {
	scoresTotal := float64(0)
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
			scoresTotal = scoresTotal + vote.Balance
			choices := []float64{}
			for _, v := range vote.Choice {
				choices = append(choices, float64(v))
			}
			for idx, value := range vote.Choice {
				choiceWeightPercent := utils.CalcPercentageOfSum(float64(value), choices)
				index, err := strconv.ParseInt(idx, 10, 64)
				if err != nil {
					log.Println("Error while parsing string:-", err)
					continue
				}
				for sIdx, score := range vote.Scores {
					choiceWeightPower := choiceWeightPercent * score
					sqrt := math.Sqrt(choiceWeightPower)
					scoresByStrategy[index-1][sIdx] = scoresByStrategy[index-1][sIdx] + sqrt
				}
			}
		}
	}

	for _, scores := range scoresByStrategy {
		for idx, score := range scores {
			scores[idx] = score * score
		}
	}

	flattenScoresByStrategy := funk.FlattenDeep(scoresByStrategy)

	for idx, scores := range scoresByStrategy {
		percentageOfScores := []float64{}
		for _, score := range scores {
			percentageOfScores = append(percentageOfScores, utils.CalcPercentageOfSum(score, flattenScoresByStrategy.([]float64)))
		}
		scoresByStrategy[idx] = utils.CalcReducedQuadraticScores(scoresTotal, percentageOfScores)
	}

	return scoresByStrategy
}
