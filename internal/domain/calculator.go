package domain

import (
	"errors"
)

func CalculateFinalScore(rubric Rubric, evaluations []Evaluation) (float64, error) {
	var totalScore float64
	seen := make(map[string]struct{})
	for _, evaluation := range evaluations {
		if _, exists := seen[evaluation.CriterionName()]; exists {
			return 0, errors.New("дубликат оценки по критерию: " + evaluation.CriterionName())
		}
		seen[evaluation.CriterionName()] = struct{}{}

		found := false
		for _, rubricCriterion := range rubric.Criteria() {
			if evaluation.CriterionName() == rubricCriterion.Name() {
				totalScore += float64(evaluation.Score()) / float64(rubricCriterion.MaxScore()) * float64(rubricCriterion.Weight())
				found = true
				break
			}
		}
		if !found {
			return 0, errors.New("оценка содержит критерий, который не существует в рубрике")
		}
	}

	return totalScore * 10, nil
}
