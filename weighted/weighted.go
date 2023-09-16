package weighted

import (
	"log"
	"math/big"
	"strconv"

	"github.com/thoas/go-funk"
)

type WeightedVote struct {
	Choice  WeightedChoice `json:"choice"`
	Balance *big.Float     `json:"balance"`
	Scores  []*big.Float   `json:"scores"`
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

func CalcPercentageOfSum(choice *big.Float, choices []*big.Float) *big.Float {
	if new(big.Float).Set(choice).SetPrec(5).Cmp(big.NewFloat(0.0).SetPrec(5)) == 0 {
		return big.NewFloat(0.0)
	}

	whole := funk.Reduce(choices, func(acc *big.Float, c *big.Float) *big.Float {
		return acc.Add(acc, new(big.Float).Set(c))
	}, big.NewFloat(0)).(*big.Float)

	if whole.SetPrec(5).Cmp(big.NewFloat(0.0).SetPrec(5)) == 0 {
		return big.NewFloat(0.0)
	}

	return new(big.Float).Set(choice).Quo(new(big.Float).Set(choice), whole)
}

func WeightedPower(choice *big.Float, choices []*big.Float, balance *big.Float) *big.Float {
	percentage := CalcPercentageOfSum(choice, choices)
	return percentage.Mul(percentage, balance)
}

func CalcReducedQuadraticScores(scoresTotal *big.Float, percentages []*big.Float) []*big.Float {
	return funk.Map(percentages, func(p *big.Float) *big.Float {
		return p.Mul(p, scoresTotal)
	}).([]*big.Float)
}

func (v *WeightedVoting) GetValidVotes() []WeightedVote {
	return funk.Filter(v.Votes, func(vote WeightedVote) bool {
		return IsValidChoice(vote.Choice, v.Choices)
	}).([]WeightedVote)
}

func (v *WeightedVoting) GetScoresTotal() *big.Float {
	return funk.Reduce(v.Votes, func(acc *big.Float, vote WeightedVote) *big.Float {
		return acc.Add(acc, vote.Balance)
	}, big.NewFloat(0)).(*big.Float)
}

func (v *WeightedVoting) GetScores() []*big.Float {
	scoresTotal := big.NewFloat(0)
	scores := []*big.Float{}

	for range v.Choices {
		scores = append(scores, big.NewFloat(0))
	}

	for _, vote := range v.Votes {
		if IsValidChoice(vote.Choice, v.Choices) {
			scoresTotal = scoresTotal.Add(scoresTotal, vote.Balance)
			choices := []*big.Float{}
			for _, v := range vote.Choice {
				choices = append(choices, big.NewFloat(float64(v)))
			}

			for idx, value := range vote.Choice {
				choiceWeightedPower := WeightedPower(big.NewFloat(float64(value)), choices, vote.Balance)
				index, err := strconv.ParseInt(idx, 10, 64)
				if err != nil {
					log.Println("Error while parsing string:-", err)
					continue
				}
				scores[index-1] = new(big.Float).Set(scores[index-1]).Add(new(big.Float).Set(scores[index-1]), choiceWeightedPower)
			}

		}
	}

	percentageOfScores := []*big.Float{}
	for _, score := range scores {
		percentageOfScores = append(percentageOfScores, CalcPercentageOfSum(new(big.Float).Set(score), scores))
	}
	newScores := CalcReducedQuadraticScores(scoresTotal, percentageOfScores)

	return newScores
}

func (v *WeightedVoting) GetScoresByStrategy() [][]*big.Float {
	scoresTotal := big.NewFloat(0)
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
			scoresTotal = scoresTotal.Add(scoresTotal, vote.Balance)
			choices := []*big.Float{}
			for _, v := range vote.Choice {
				choices = append(choices, big.NewFloat(float64(v)))
			}
			for idx, value := range vote.Choice {
				index, err := strconv.ParseInt(idx, 10, 64)
				if err != nil {
					log.Println("Error while parsing string:-", err)
					continue
				}
				for sIdx, score := range vote.Scores {
					choiceWeightedPower := WeightedPower(big.NewFloat(float64(value)), choices, score)
					scoresByStrategy[index-1][sIdx] = scoresByStrategy[index-1][sIdx].Add(scoresByStrategy[index-1][sIdx], choiceWeightedPower)
				}
			}
		}
	}

	flattenScoresByStrategy := funk.FlattenDeep(scoresByStrategy)

	for idx, scores := range scoresByStrategy {
		percentageOfScores := []*big.Float{}
		for _, score := range scores {
			percentageOfScores = append(percentageOfScores, CalcPercentageOfSum(new(big.Float).Set(score), flattenScoresByStrategy.([]*big.Float)))
		}
		scoresByStrategy[idx] = CalcReducedQuadraticScores(scoresTotal, percentageOfScores)
	}

	return scoresByStrategy
}
