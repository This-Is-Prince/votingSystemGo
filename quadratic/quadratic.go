package quadratic

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/thoas/go-funk"
)

type QuadraticVote struct {
	Choice  QuadraticChoice `json:"choice"`
	Balance *big.Float      `json:"balance"`
	Scores  []*big.Float    `json:"scores"`
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

func CalcPercentageOfSum(choice *big.Float, choices []*big.Float) *big.Float {
	if new(big.Float).Set(choice).SetPrec(5).Cmp(big.NewFloat(0.0).SetPrec(5)) == 0 {
		return big.NewFloat(0.0)
	}

	whole := funk.Reduce(choices, func(acc *big.Float, c *big.Float) *big.Float {
		return acc.Add(acc, c)
	}, big.NewFloat(0)).(*big.Float)

	if whole.SetPrec(5).Cmp(big.NewFloat(0.0).SetPrec(5)) == 0 {
		return big.NewFloat(0.0)
	}

	return choice.Quo(choice, whole)
}

func CalcReducedQuadraticScores(scoresTotal *big.Float, percentages []*big.Float) []*big.Float {
	return funk.Map(percentages, func(p *big.Float) *big.Float {
		return p.Mul(p, scoresTotal)
	}).([]*big.Float)
}

func (v *QuadraticVoting) GetValidVotes() []QuadraticVote {
	return funk.Filter(v.Votes, func(vote QuadraticVote) bool {
		return IsValidChoice(vote.Choice, v.Choices)
	}).([]QuadraticVote)
}

func (v *QuadraticVoting) GetScoresTotal() *big.Float {
	return funk.Reduce(v.Votes, func(acc *big.Float, vote QuadraticVote) *big.Float {
		return acc.Add(acc, vote.Balance)
	}, big.NewFloat(0)).(*big.Float)
}

func (v *QuadraticVoting) GetScores() []*big.Float {
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
				choiceWeightPercent := CalcPercentageOfSum(big.NewFloat(float64(value)), choices)

				sqrt := big.NewFloat(0).Sqrt(big.NewFloat(0).Mul(choiceWeightPercent, vote.Balance))
				index, err := strconv.ParseInt(idx, 10, 64)
				if err != nil {
					fmt.Println("Error while parsing string:-", err)
					continue
				}
				scores[index-1] = scores[index-1].Add(scores[index-1], sqrt)
			}
		}
	}

	for idx, score := range scores {
		scores[idx] = score.Mul(score, score)
	}

	percentageOfScores := []*big.Float{}
	for _, score := range scores {
		percentageOfScores = append(percentageOfScores, CalcPercentageOfSum(new(big.Float).Set(score), scores))
	}
	return CalcReducedQuadraticScores(scoresTotal, percentageOfScores)
}

func (v *QuadraticVoting) GetScoresByStrategy() [][]*big.Float {
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
				choiceWeightPercent := CalcPercentageOfSum(big.NewFloat(float64(value)), choices)
				index, err := strconv.ParseInt(idx, 10, 64)
				if err != nil {
					fmt.Println("Error while parsing string:-", err)
					continue
				}
				for sIdx, score := range vote.Scores {
					sqrt := big.NewFloat(0).Sqrt(big.NewFloat(0).Mul(choiceWeightPercent, score))
					scoresByStrategy[index-1][sIdx] = scoresByStrategy[index-1][sIdx].Add(scoresByStrategy[index-1][sIdx], sqrt)
				}
			}
		}
	}

	for _, scores := range scoresByStrategy {
		for idx, score := range scores {
			scores[idx] = score.Mul(score, score)
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
